package languages

import (
	"fmt"
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/shell"
)

type CppTools struct {
	Shell *shell.ShellHandler

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
		tool.installFunc = func() error {
			fmt.Println("Installing ninja...")
			err := c.Shell.Brew("ninja")
			if err != nil {
				return fmt.Errorf("failed to install ninja: %w", err)
			}
			fmt.Println("ninja installed successfully")
			return nil
		}
		c.Ninja = tool
	case "cmake":
		tool.installFunc = func() error {
			fmt.Println("Installing cmake...")
			err := c.Shell.Brew("cmake")
			if err != nil {
				return fmt.Errorf("failed to install cmake: %w", err)
			}
			fmt.Println("cmake installed successfully")
			return nil
		}
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
