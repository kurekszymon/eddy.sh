package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/types"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

var DebugEnabled = strings.EqualFold(os.Getenv("EDDY_DEBUG"), "1")

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			utils.Log("Failed to get user home directory: "+err.Error(), types.LogError)
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
