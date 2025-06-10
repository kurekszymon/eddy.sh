//go:build windows
// +build windows

package shell

import "errors"

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
