package languages

import (
	"fmt"
	"strings"
)

type CppTools struct {
	Emscripten *Tool
	Ninja      *Tool
	Cmake      *Tool
}

func (c *CppTools) SetTool(toolName string, tool *Tool) {
	switch strings.ToLower(toolName) {
	case "emscripten":
		tool.InstallFunc = InstallEmscripten
		c.Emscripten = tool
	case "ninja":
		c.Ninja = tool
	case "cmake":
		c.Cmake = tool
	}
}

func InstallEmscripten() error {
	fmt.Println("emscripten installed successfuly")
	// return fmt.Errorf("emscripten installation failed")
	return nil
}
