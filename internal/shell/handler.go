package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/kurekszymon/eddy.sh/internal/exit_codes"
	"github.com/kurekszymon/eddy.sh/internal/globals"
	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

func NewShellHandler() Shell {
	return &ShellHandler{}
}

func (s *ShellHandler) GetEddyDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	eddyDir := utils.ExpandPath(filepath.Join(homeDir, ".eddy.sh"))

	err = utils.EnsureDir(eddyDir)
	if err != nil {
		return "", err
	}

	return eddyDir, nil
}

func (s *ShellHandler) GetEddyBinDir() (string, error) {
	eddyDir, err := s.GetEddyDir()
	if err != nil {
		return "", err
	}

	eddyBinDir := utils.ExpandPath(filepath.Join(eddyDir, "bin"))

	err = utils.EnsureDir(eddyBinDir)
	if err != nil {
		return "", err
	}

	return eddyBinDir, nil
}

func (s *ShellHandler) GitClone(repoURL string, cloneDir string) error {
	repoName := filepath.Base(repoURL)
	if filepath.Ext(repoName) == ".git" {
		repoName = repoName[:len(repoName)-len(".git")]
	}

	cloneDir = utils.ExpandPath(filepath.Join(cloneDir, repoName))

	err := s.run("git clone --progress %s %s 2>&1", repoURL, utils.ExpandPath(cloneDir))
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Successfully cloned %s into %s\n", repoURL, cloneDir)
	logger.Info(message)
	return nil
}

func (s *ShellHandler) Curl(url string, outputDir string) error {
	command := fmt.Sprintf("curl -fL --progress-bar --output-dir %s -O %s", outputDir, url)
	err := s.run(command)
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

func (s *ShellHandler) Unzip(archivePath string, targetDir string) error {
	err := utils.EnsureDir(targetDir)
	if err != nil {
		return err
	}

	err = s.run("tar -xf %s -C %s", archivePath, targetDir)

	if err != nil {
		return err
	}

	return nil
}

func (s *ShellHandler) RunScriptFile(filename string) error {
	filename = utils.ExpandPath(filename)
	os.Chmod(filename, 0755)

	err := s.run(filename)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShellHandler) RunCustomScript(script string) error {
	message := fmt.Sprintf("You are about to run custom command: %s", script)
	logger.Warn(message)
	utils.PromptConfirm("Custom script can potentially harm your system. Do you want to continue? (Y/N) ",
		"ERROR: User aborted running custom script",
		exit_codes.CUSTOM_SCRIPT_EXIT)

	err := s.run(script)
	if err != nil {
		return fmt.Errorf("failed to run custom script %s: %w", script, err)
	}

	return nil
}

func (s *ShellHandler) RunScriptFileInDir(filename string, dir string, args ...string) error {
	scriptPath := utils.ExpandPath(filepath.Join(dir, filename))

	os.Chmod(scriptPath, 0755)

	command := scriptPath
	for _, arg := range args {
		// hack to run script with all the args passed to it
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
	command = utils.FormatArgs(command, args)

	if globals.DebugEnabled {
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

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command '%s': %w", command, err)
	}

	s.handlePipes(stdout, stderr)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command '%s' failed: %w", command, err)
	}

	return nil
}

func (s *ShellHandler) handlePipes(stdout io.Reader, stderr io.Reader) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, stdout)
	}()

	go func() {
		defer wg.Done()
		io.Copy(os.Stderr, stderr)
	}()

	wg.Wait()
}
