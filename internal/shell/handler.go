package shell

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type ShellHandler struct {
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
