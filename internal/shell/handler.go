package shell

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type ShellHandler struct {
	Verbose bool // If true, prints commands before executing
}

// Run executes a shell command and returns its output or an error
func (s *ShellHandler) Run(command string, args ...string) error {
	if s.Verbose {
		fmt.Printf("Executing: %s %s\n", command, args)
	}

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
	s.HandlePipes(stdout, stderr)
	cmd.Wait()
	return nil
}

func (s *ShellHandler) HandlePipes(stdout io.Reader, stderr io.Reader) {
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
