//go:build windows
// +build windows

package shell

type ShellHandler struct{}

// TODO: Example commands, need to be adapted for Windows

func (s *ShellHandler) CheckCommand(command string) error {
	// Windows does not have a direct equivalent to 'command -v'
	// We can use 'where' to check if a command exists
	err := s.run("where %s > NUL", command)

	if err != nil {
		return err
	}

	return nil
}

// add echo method
func (s *ShellHandler) Echo(message string) error {
	// Windows does not have a direct equivalent to 'echo'
	// We can use 'echo' command to print messages
	dir, err := s.GetEddyDir()
	if err != nil {
		return err
	}
	err = s.run("echo %s %s", message, dir)

	if err != nil {
		return err
	}

	return nil
}
