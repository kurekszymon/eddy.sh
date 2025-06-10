package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

	err := s.run("git clone %s %s", repoURL, ExpandPath(cloneDir))
	if err != nil {
		return fmt.Errorf("failed to clone repository %s into %s: %w", repoURL, cloneDir, err)
	}

	fmt.Printf("Successfully cloned %s into %s\n", repoURL, cloneDir)
	return nil
}

func (s *ShellHandler) Curl(url string) error {
	eddyDir, err := s.GetEddyDir()
	if err != nil {
		return fmt.Errorf("failed to get eddy dir: %w", err)
	}

	command := fmt.Sprintf("curl -fL --output-dir %s -O %s", eddyDir, url)
	err = s.run(command)
	if err != nil {
		return fmt.Errorf("failed to run curl: %w", err)
	}
	return nil
}

func (s *ShellHandler) Echo(message string) error {
	err := s.run("echo %s", message)
	if err != nil {
		return fmt.Errorf("failed to echo %s: %w", message, err)
	}
	return nil
}

func (s *ShellHandler) Unzip(filename string, target_dir string) error {
	eddyDir, err := s.GetEddyDir()
	if err != nil {
		return fmt.Errorf("failed to get eddy dir: %w", err)
	}

	filename = filepath.Join(eddyDir, filename)
	target_dir = filepath.Join(eddyDir, target_dir)
	err = s.ensureDir(target_dir)
	if err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", target_dir, err)
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
		return fmt.Errorf("failed to run script file: %w", err)
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
		fmt.Println("DEBUG: Running command:", command)
	}

	cmd := exec.Command("sh", "-c", command)
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
					fmt.Printf("Error reading stdout: %v\n", err)
				}
				close(stdoutChan)
				return
			}

			if DebugEnabled {
				line = "DEBUG stdout: " + line
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
					fmt.Printf("Error reading stderr: %v\n", err)
				}
				close(stderrChan)
				return
			}
			if DebugEnabled {
				line = "DEBUG: stderr: " + line
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
