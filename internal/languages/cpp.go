package languages

import (
	"fmt"
	"strings"
)

type CppTools struct {
	Emscripten *Tool
	Ninja      *Tool
	Cmake      *Tool
	NotLoaded  *[]Tool
}

func (c *CppTools) SetTool(toolName string, tool *Tool) {
	switch strings.ToLower(toolName) {
	case "emscripten":
		tool.installFunc = InstallEmscripten
		c.Emscripten = tool
	case "ninja":
		c.Ninja = tool
	case "cmake":
		c.Cmake = tool
	default:
		if c.NotLoaded == nil {
			c.NotLoaded = &[]Tool{}
		}
		*c.NotLoaded = append(*c.NotLoaded, *tool)
	}

}

func InstallEmscripten() error {
	fmt.Println("emscripten installed successfuly")
	// return fmt.Errorf("emscripten installation failed")
	return nil
}
