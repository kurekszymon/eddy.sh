package languages

import "fmt"

type Tool struct {
	Name        string
	Version     string
	InstallFunc func() error
}

func (t *Tool) Install() error {
	if t.InstallFunc != nil {
		return t.InstallFunc()
	}
	return fmt.Errorf("Install not implemented for tool: %s", t.Name)
}

type ToolSetter interface {
	SetTool(toolName string, tool *Tool)
}
