package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

var DebugEnabled = strings.EqualFold(os.Getenv("EDDY_DEBUG"), "1")
