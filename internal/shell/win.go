//go:build windows
// +build windows

package shell

import (
	"errors"
	"path/filepath"
)

type ShellHandler struct{}

func (s *ShellHandler) CheckCommand(command string) error {
	err := s.run("where %s > NUL", command)

	if err != nil {
		return err
	}

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

	s.run("if not exist %s", link_dir, "mklink", link_dir, source)

	return nil
}
