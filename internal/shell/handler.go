package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/kurekszymon/eddy.sh/internal/error_codes"
	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

func (s *ShellHandler) GetEddyDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	eddyDir := ExpandPath(filepath.Join(homeDir, ".eddy.sh"))

	err = s.ensureDir(eddyDir)
	if err != nil {
		return "", err
	}

	return eddyDir, nil
}

func (s *ShellHandler) GitClone(repoURL string, cloneDir string) error {
	repoName := filepath.Base(repoURL)
	if filepath.Ext(repoName) == ".git" {
		repoName = repoName[:len(repoName)-len(".git")]
	}

	cloneDir = ExpandPath(filepath.Join(cloneDir, repoName))

	err := s.run("git clone --progress %s %s 2>&1", repoURL, ExpandPath(cloneDir))
	if err != nil {
		return fmt.Errorf("failed to clone repository %s into %s: %w", repoURL, cloneDir, err)
	}

	message := fmt.Sprintf("Successfully cloned %s into %s\n", repoURL, cloneDir)
	logger.Info(message)
	return nil
}

func (s *ShellHandler) Curl(url string) error {
	eddyDir, err := s.GetEddyDir()
	if err != nil {
		return err
	}

	command := fmt.Sprintf("curl -fL --output-dir %s -O %s", eddyDir, url)
	err = s.run(command)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShellHandler) Echo(message string) error {
	err := s.run("echo %s", message)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShellHandler) Unzip(filename string, target_dir string) error {
	eddyDir, err := s.GetEddyDir()
	if err != nil {
		return err
	}

	filename = filepath.Join(eddyDir, filename)
	target_dir = filepath.Join(eddyDir, target_dir)
	err = s.ensureDir(target_dir)
	if err != nil {
		return err
	}

	err = s.run("tar -xf %s -C %s", filename, target_dir)

	if err != nil {
		return err
	}

	return nil
}

func (s *ShellHandler) RunScriptFile(filename string) error {
	filename = ExpandPath(filename)
	os.Chmod(filename, 0755)

	err := s.run(filename)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShellHandler) RunCustomScript(script string) error {
	// potentially flag it in config / omit custom script checks

	message := fmt.Sprintf("-- You are about to run custom command: %s \n", script)
	logger.Warn(message)
	utils.PromptConfirm("-- Custom script can potentially harm your system. Do you want to continue? (Y/N) ",
		"-- ERROR: User aborted running custom script",
		error_codes.CUSTOM_SCRIPT_EXIT)

	err := s.run(script)
	if err != nil {
		return fmt.Errorf("-- failed to run custom script %s: %w", script, err)
	}

	return nil
}

func (s *ShellHandler) RunScriptFileInDir(filename string, dir string, args ...string) error {
	scriptPath := ExpandPath(filepath.Join(dir, filename))

	os.Chmod(scriptPath, 0755)

	// hack to run script with all the args passed to it
	command := scriptPath
	for _, arg := range args {
		command = fmt.Sprintf("%s %s", command, arg)
	}

	err := s.run(command)
	if err != nil {
		return fmt.Errorf("failed to run script file in directory: %w", err)
	}
	return nil
}

// This is the main function to run a command in the shell.
// Note: The command is executed in a shell, so shell features like pipes and redirection are available
func (s *ShellHandler) run(command string, args ...string) error {
	command = FormatArgs(command, args)

	if DebugEnabled {
		logger.Debug("Running command: " + command)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Env = append(os.Environ(), "PYTHONUNBUFFERED=1")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	cmd.Start()
	s.handlePipes(stdout, stderr)
	err = cmd.Wait()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() != 0 {
				return fmt.Errorf("command %s failed with exit code %d", command, exitError.ExitCode())
			}
		}
	}

	return nil
}

func (s *ShellHandler) ensureDir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

func (s *ShellHandler) handlePipes(stdout io.Reader, stderr io.Reader) {
	stdoutChan := make(chan string)
	stderrChan := make(chan string)

	go func() {
		stdoutReader := bufio.NewReader(stdout)
		for {
			line, err := stdoutReader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					logger.Error("Error reading from stdout")
				}
				close(stdoutChan)
				return
			}

			if DebugEnabled {
				line = logger.FormatLogType(line, logger.LogDebug)
			}
			stdoutChan <- line
		}
	}()

	go func() {
		stderrReader := bufio.NewReader(stderr)
		for {
			line, err := stderrReader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					logger.Info("Error reading from stderr")
				}
				close(stderrChan)
				return
			}
			if DebugEnabled {
				line = logger.FormatLogType(line, logger.LogDebug)
			}
			stderrChan <- line
		}
	}()

	for {
		select {
		case line, ok := <-stdoutChan:
			if !ok {
				stdoutChan = nil
			} else {
				fmt.Print(line)
			}
		case line, ok := <-stderrChan:
			if !ok {
				stderrChan = nil
			} else {
				fmt.Print(line)
			}
		}

		if stdoutChan == nil && stderrChan == nil {
			break
		}
	}
}
