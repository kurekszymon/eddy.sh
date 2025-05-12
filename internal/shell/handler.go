package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

type ShellHandler struct{}

// using code down there handle getting eddy.sh dir that should resolve to ~/.eddy.sh
func (s *ShellHandler) GetEddyDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	eddyDir := filepath.Join(homeDir, ".eddy.sh")
	// check if the directory exists and create it if it doesn't
	// you can use s.existsDir() method here
	s.ensureDir(eddyDir)
	if err != nil {
		return "", fmt.Errorf("failed to create eddy dir: %w", err)
	}

	return eddyDir, nil
}

func (s *ShellHandler) RunScriptFile(filename string) error {
	s.run("chmod +x $0", filename)
	err := s.run(filename)
	if err != nil {
		return fmt.Errorf("failed to run script file: %w", err)
	}
	return nil
}

func (s *ShellHandler) Curl(url string) error {
	eddyDir, err := s.GetEddyDir()
	if err != nil {
		return fmt.Errorf("failed to get eddy dir: %w", err)
	}
	err = s.run("curl -L --output-dir $0 -O $1", eddyDir, "https://github.com/ninja-build/ninja/releases/download/v1.12.1/ninja-mac.zip")
	if err != nil {
		return fmt.Errorf("failed to run curl: %w", err)
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

	err = s.run("tar -xf $0 -C $1", filename, target_dir)
	if err != nil {
		return fmt.Errorf("failed to run tar -xf: %w", err)
	}
	return nil
}

func (s *ShellHandler) Echo(message string) error {
	err := s.run("echo $0", message)
	if err != nil {
		return fmt.Errorf("failed to echo %s: %w", message, err)
	}
	return nil
}

func (s *ShellHandler) Brew(pkg string) error {
	err := s.run("brew install $0", pkg)
	if err != nil {
		return fmt.Errorf("failed to install %s: %w", pkg, err)
	}
	return nil
}

// This is the main function to run a command in the shell.
// Note: The command is executed in a shell, so shell features like pipes and redirection are available
// but the command should be passed as a single string
// For example, to run "echo hello {name}", you would call:
// Run("echo hello $0", "{name}")
func (s *ShellHandler) run(command string, args ...string) error {

	fullCommand := append([]string{"-c", command}, args...)
	cmd := exec.Command("sh", fullCommand...)
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
	cmd.Wait()
	return nil
}

func (s *ShellHandler) ensureDir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
		return nil
	}
	return fmt.Errorf("failed to check directory %s: %w", path, err)
}

func (s *ShellHandler) handlePipes(stdout io.Reader, stderr io.Reader) {
	stdout_reader := bufio.NewReader(stdout)
	line, err := stdout_reader.ReadString('\n')

	for err == nil {
		line = "stdout: " + line
		fmt.Print(line)
		line, err = stdout_reader.ReadString('\n')
	}

	stderr_reader := bufio.NewReader(stderr)
	er_line, er_err := stderr_reader.ReadString('\n')
	for er_err == nil {
		er_line = "stderr: " + er_line
		fmt.Print(er_line)
		er_line, er_err = stderr_reader.ReadString('\n')
	}
}
