package languages

import (
	"strings"
)

type JsTools struct {
	Nvm       *Tool
	NotLoaded *[]Tool
}

func (j *JsTools) SetTool(toolName string, tool *Tool) {
	switch strings.ToLower(toolName) {
	case "nvm":
		j.Nvm = tool
	default:
		if j.NotLoaded == nil {
			j.NotLoaded = &[]Tool{}
		}
		*j.NotLoaded = append(*j.NotLoaded, *tool)
	}
}
