package javascript

import (
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/shell"
)

type Tools struct {
	Shell *shell.ShellHandler

	Nvm       *installers.Tool
	NotLoaded *[]installers.Tool
}

func (j *Tools) SetTool(toolName string, tool *installers.Tool) {
	switch strings.ToLower(toolName) {
	case "nvm":
		j.Nvm = tool
	default:
		if j.NotLoaded == nil {
			j.NotLoaded = &[]installers.Tool{}
		}
		*j.NotLoaded = append(*j.NotLoaded, *tool)
	}
}
