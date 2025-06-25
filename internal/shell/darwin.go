//go:build darwin
// +build darwin

package shell

import (
	"fmt"
	"path/filepath"

	"github.com/kurekszymon/eddy.sh/internal/logger"
)

type ShellHandler struct{}

func (s *ShellHandler) CheckCommand(command string) error {
	err := s.run("command -v %s > /dev/null", command)

	if err != nil {
		return fmt.Errorf("command %s not found: %w", command, err)
	}

	return nil
}

func (s *ShellHandler) Brew(pkg string) error {
	err := s.run("brew install %s", pkg)
	if err != nil {
		return fmt.Errorf("failed to install %s: %w", pkg, err)
	}
	return nil
}

func (s *ShellHandler) Symlink(source string, dest string) error {
	eddy_dir, err := s.GetEddyDir()

	if err != nil {
		return err
	}

	eddy_bin := filepath.Join(eddy_dir, "bin")
	err = s.ensureDir(eddy_bin)
	if err != nil {
		return err
	}

	link_dir := filepath.Join(eddy_bin, dest)

	logger.Info("Creating symlink for " + dest)
	s.run("ln -s %s %s", source, link_dir)
	return nil
}
