package cpp

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/kurekszymon/eddy.sh/internal/types"
)

func (c *Tools) EmscriptenInstall() error {
	if c.PkgManager == types.Brew {
		return c.brewEmsdk()
	}

	return c.manualEmsdk()
}

func (c *Tools) brewEmsdk() error {
	fmt.Println("Installing emscripten with brew...")
	err := c.Shell.Brew("emscripten")
	if err != nil {
		return err
	}
	fmt.Println("emscripten installed successfully")
	return nil
}

func (c *Tools) manualEmsdk() error {
	err := c.Shell.CheckCommand("git")

	if err != nil {
		return errors.New("git is not installed, please install git to proceed with emscripten installation")
	}

	fmt.Println("Cloning emscripten repository...")
	err = c.Shell.GitClone("https://github.com/emscripten-core/emsdk.git", c.CloneDir)

	if err != nil {
		return err
	}

	c.CloneDir = filepath.Join(c.CloneDir, "emsdk")

	fmt.Println("Running emscripten install script...")
	err = c.Shell.RunScriptFileInDir("emsdk", c.CloneDir, "install", "latest")

	if err != nil {
		return err
	}

	fmt.Println("Activating emscripten environment...")
	err = c.Shell.RunScriptFileInDir("emsdk", c.CloneDir, "activate", "latest")

	if err != nil {
		return err
	}

	fmt.Println("SUCCESS: Emscripten installed successfully")

	return nil
}
