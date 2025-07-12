package javascript

import (
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

type Tools struct {
	Shell      *shell.ShellHandler
	PkgManager types.PkgManager

	Nvm       *installers.Tool
	NotLoaded *[]installers.Tool
}

func (j *Tools) SetTool(toolName string, tool *installers.Tool) {
	switch strings.ToLower(toolName) {
	case "nvm":
		tool.InstallFunc = j.NvmInstall
		j.Nvm = tool
	default:
		if j.NotLoaded == nil {
			j.NotLoaded = &[]installers.Tool{}
		}
		*j.NotLoaded = append(*j.NotLoaded, *tool)
	}
}

func (c *Tools) Install() map[string]error {
	errors := make(map[string]error)

	if c.Nvm != nil {
		if err := c.Nvm.Install(); err != nil {
			errors["Nvm"] = err
		}
	}

	return errors
}
