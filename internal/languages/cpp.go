package languages

import (
	"fmt"
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/globals"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

type CppTools struct {
	Shell    *shell.ShellHandler
	Platform types.Platform

	Emscripten *Tool
	Ninja      *Tool
	Cmake      *Tool
	NotLoaded  *[]Tool
}

func (c *CppTools) SetTool(toolName string, tool *Tool) {
	switch strings.ToLower(toolName) {
	case "emscripten":
		tool.installFunc = c.emscriptenInstall
		c.Emscripten = tool
	case "ninja":
		tool.installFunc = c.ninjaInstall
		c.Ninja = tool
	case "cmake":
		tool.installFunc = c.cmakeInstall
		c.Cmake = tool

	default:
		if c.NotLoaded == nil {
			c.NotLoaded = &[]Tool{}
		}
		*c.NotLoaded = append(*c.NotLoaded, *tool)
	}

}

func (c *CppTools) cmakeInstall() error {
	if c.Platform == types.Brew {
		fmt.Println("Installing cmake...")
		err := c.Shell.Brew("cmake")
		if err != nil {
			return fmt.Errorf("failed to install cmake: %w", err)
		}
		fmt.Println("cmake installed successfully")
		return nil
	}

	fmt.Println("Manual installation of cmake is not supported yet.")
	fmt.Println("Please follow the instructions at https://cmake.org/install/")
	return nil
}

func (c *CppTools) ninjaInstall() error {
	if c.Platform == types.Brew {
		fmt.Println("Installing ninja...")
		err := c.Shell.Brew("ninja")
		if err != nil {
			return fmt.Errorf("failed to install ninja: %w", err)
		}
		fmt.Println("ninja installed successfully")
		return nil
	}

	ninja_url := fmt.Sprintf("https://github.com/ninja-build/ninja/releases/tag/v%s/%s.zip", c.Ninja.Version, globals.NINJA_DIRNAME)
	err := c.Shell.Curl(ninja_url)
	if err != nil {
		return fmt.Errorf("failed to download ninja: %w", err)
	}

	err = c.Shell.Unzip(globals.NINJA_DIRNAME, "ninja")

	if err != nil {
		return fmt.Errorf("failed to unzip ninja: %w", err)
	}

	return nil
}

func (c *CppTools) emscriptenInstall() error {
	if c.Platform == types.Brew {
		fmt.Println("Installing emscripten...")
		err := c.Shell.Brew("emscripten")
		if err != nil {
			return fmt.Errorf("failed to install emscripten: %w", err)
		}
		fmt.Println("emscripten installed successfully")
		return nil
	}
	fmt.Println("Manual installation of emscripten is not supported yet.")
	fmt.Println("Please follow the instructions at https://emscripten.org/docs/getting_started/downloads.html")
	return nil
}
