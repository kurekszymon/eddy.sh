package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	err := s.run("git clone $0 $1", repoURL, ExpandPath(cloneDir))
	if err != nil {
		return fmt.Errorf("failed to clone repository %s into %s: %w", repoURL, cloneDir, err)
	}

	fmt.Printf("Successfully cloned %s into %s\n", repoURL, cloneDir)
	return nil
}

// This is the main function to run a command in the shell.
// Note: The command is executed in a shell, so shell features like pipes and redirection are available
// but the command should be passed as a single string
// For example, to run "echo hello {name}", you would call:
// Run("echo hello $0", "{name}")
func (s *ShellHandler) run(command string, args ...string) error {
	fullCommand := append([]string{"-c", command}, args...)

	if DebugEnabled {
		fmt.Println("DEBUG: Running command:", strings.Join(fullCommand, " "))
	}
	cmd := exec.Command("sh", fullCommand...)
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
