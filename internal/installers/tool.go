package installers

import "fmt"

type Tool struct {
	Name        string
	Version     string
	installFunc func() error
}

func (t *Tool) Install() error {
	if t.installFunc != nil {
		return t.installFunc()
	}
	return fmt.Errorf("Install not implemented for tool: %s", t.Name)
}

type ToolSetter interface {
	SetTool(toolName string, tool *Tool)
}
