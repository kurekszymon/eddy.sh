package javascript

import (
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

type Installer struct {
	Shell      shell.Shell
	PkgManager types.PkgManager

	Available map[string]*installers.Tool
	NotLoaded []installers.Tool
}

func (j *Installer) SetTool(toolName string, tool *installers.Tool) {
	toolKey := strings.ToLower(toolName)

	if j.Available == nil {
		j.Available = make(map[string]*installers.Tool)
	}

	switch toolKey {
	case "nvm":
		tool.InstallFunc = j.NvmInstall
	default:
		j.NotLoaded = append(j.NotLoaded, *tool)
		return
	}

	j.Available[toolKey] = tool
}

func (c *Installer) Install() map[string]error {
	errors := make(map[string]error)

	for name, tool := range c.Available {
		if err := tool.Install(); err != nil {
			errors[name] = err
		}
	}

	return errors
}
