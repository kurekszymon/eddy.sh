package cpp

import (
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

type Installer struct {
	Shell      *shell.ShellHandler
	PkgManager types.PkgManager

	Available map[string]*installers.Tool
	NotLoaded []installers.Tool
}

func (c *Installer) SetTool(toolName string, tool *installers.Tool) {
	toolKey := strings.ToLower(toolName)

	if c.Available == nil {
		c.Available = make(map[string]*installers.Tool)
	}

	switch toolKey {
	case "emscripten":
		tool.InstallFunc = c.EmscriptenInstall
	case "ninja":
		tool.InstallFunc = c.NinjaInstall
	case "cmake":
		tool.InstallFunc = c.CmakeInstall
	default:
		c.NotLoaded = append(c.NotLoaded, *tool)
		return
	}

	c.Available[toolKey] = tool
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
