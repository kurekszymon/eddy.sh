package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var DebugEnabled = strings.EqualFold(os.Getenv("EDDY_DEBUG"), "1")

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func FormatArgs(command string, args []string) string {
	interfaceArgs := make([]any, len(args))
	for i, v := range args {
		interfaceArgs[i] = v
	}

	return fmt.Sprintf(command, interfaceArgs...)
}
