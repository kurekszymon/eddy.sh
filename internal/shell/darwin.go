//go:build darwin
// +build darwin

package shell

import (
	"fmt"
	"path/filepath"
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
