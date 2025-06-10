//go:build darwin
// +build darwin

package shell

import (
	"fmt"
)

type ShellHandler struct{}

func (s *ShellHandler) CheckCommand(command string) error {
	// probably won't work on windows. replace $0 with %VARIABLE% and NUL.
	// educate yourself on how to do this
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
