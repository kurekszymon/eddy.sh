package cpp

import (
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

type Tools struct {
	Shell      *shell.ShellHandler
	PkgManager types.PkgManager
	CloneDir   string

	Emscripten *installers.Tool
	Ninja      *installers.Tool
	Cmake      *installers.Tool
	NotLoaded  *[]installers.Tool
}

func (c *Tools) SetTool(toolName string, tool *installers.Tool) {
	switch strings.ToLower(toolName) {
	case "emscripten":
		tool.InstallFunc = c.EmscriptenInstall
		c.Emscripten = tool
	case "ninja":
		tool.InstallFunc = c.NinjaInstall
		c.Ninja = tool
	case "cmake":
		tool.InstallFunc = c.CmakeInstall
		c.Cmake = tool

	default:
		if c.NotLoaded == nil {
			c.NotLoaded = &[]installers.Tool{}
		}
		*c.NotLoaded = append(*c.NotLoaded, *tool)
	}

}

func (c *Tools) Install() map[string]error {
	errors := make(map[string]error)

	if c.Emscripten != nil {
		if err := c.Emscripten.Install(); err != nil {
			errors["Emscripten"] = err
		}
	}

	if c.Ninja != nil {
		if err := c.Ninja.Install(); err != nil {
			errors["Ninja"] = err
		}
	}

	if c.Cmake != nil {
		if err := c.Cmake.Install(); err != nil {
			errors["CMake"] = err
		}
	}

	return errors
}
