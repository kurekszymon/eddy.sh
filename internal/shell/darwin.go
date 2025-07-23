//go:build darwin
// +build darwin

package shell

import (
	"fmt"

	"github.com/kurekszymon/eddy.sh/internal/logger"
)

type ShellHandler struct{}

func (s *ShellHandler) CheckCommand(command string) error {
	msg := fmt.Sprintf("checking for '%s'", command)
	logger.Info(msg)

	err := s.run("command -v %s > /dev/null", command)

	if err != nil {
		return fmt.Errorf("command %s not found", command)
	}

	msg = fmt.Sprintf("'%s' is present", command)
	logger.Info(msg)
	return nil
}

func (s *ShellHandler) Brew(pkg string) error {
	err := s.run("brew install %s", pkg)
	if err != nil {
		return fmt.Errorf("failed to install %s: %w", pkg, err)
	}
	return nil
}

func (s *ShellHandler) Symlink(source string, linkPath string) error {
	logger.Info("Creating symlink for " + source)
	err := s.run("ln -s %s %s", source, linkPath)
	return err
}
