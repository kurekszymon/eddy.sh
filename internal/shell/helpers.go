package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/logger"
)

var DebugEnabled = strings.EqualFold(os.Getenv("EDDY_DEBUG"), "1")

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Error("Failed to get user home directory: " + err.Error())
			return path
		}
		return filepath.Join(home, path[2:])
	}
	abs, err := filepath.Abs(path)
	if err == nil {
		return abs
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
