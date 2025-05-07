package languages

import (
	"strings"
)

type JsTools struct {
	Nvm *Tool
}

func (j *JsTools) SetTool(toolName string, tool *Tool) {
	switch strings.ToLower(toolName) {
	case "nvm":
		j.Nvm = tool
	}
}
