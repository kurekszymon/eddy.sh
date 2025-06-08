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

type ShellHandler struct{}

func (s *ShellHandler) CheckCommand(command string) error {
	// probably won't work on windows. replace $0 with %VARIABLE% and NUL.
	// educate yourself on how to do this
	err := s.run("command -v $0 > /dev/null", command)

	if err != nil {
		return fmt.Errorf("command %s not found: %w", command, err)
	}

	return nil
}

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

func (s *ShellHandler) ChmodX(filename string) error {
	fmt.Println("running chmod +x on", ExpandPath(filename))
	err := s.run("chmod +x $0", ExpandPath(filename))
	if err != nil {
		return fmt.Errorf("failed to make file executable: %w", err)
	}
	return nil
}

func (s *ShellHandler) RunScriptFile(filename string) error {
	filename = ExpandPath(filename)
	s.ChmodX(filename)

	err := s.run(filename)
	if err != nil {
		return fmt.Errorf("failed to run script file: %w", err)
	}
	return nil
}

func (s *ShellHandler) RunScriptFileInDir(filename string, dir string, args ...string) error {
	scriptPath := ExpandPath(filepath.Join(dir, filename))
	args = append([]string{scriptPath}, args...)

	s.ChmodX(scriptPath)

	command := fmt.Sprintf("%s $@", scriptPath)
	err := s.run(command, args...)
	if err != nil {
		return fmt.Errorf("failed to run script file in directory: %w", err)
	}
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
		return err
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

	if os.Getenv("EDDY_DEBUG") == "1" {
		fmt.Println("Running command:", strings.Join(fullCommand, " "))
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

			if os.Getenv("EDDY_DEBUG") == "1" {
				line = "stderr: " + line
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
			if os.Getenv("EDDY_DEBUG") == "1" {
				line = "stderr: " + line
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
