//go:build windows
// +build windows

package shell

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/kurekszymon/eddy.sh/internal/logger"
)

type ShellHandler struct{}

func (s *ShellHandler) CheckCommand(command string) error {
	msg := fmt.Sprintf("checking for '%s'", command)
	logger.Info(msg)

	err := s.run("where %s > NUL", command)

	if err != nil {
		return err
	}

	msg = fmt.Sprintf("'%s' is present", command)
	logger.Info(msg)

	return nil
}

func (s *ShellHandler) Brew(pkg string) error {
	return errors.New("brew is not supported on Windows, please use a different package manager")
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

	s.run("if not exist %s mklink %s %s", link_dir, link_dir, source)

	return nil
}
