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
		return fmt.Errorf("failed to install emscripten: %w", err)
	}
	fmt.Println("emscripten installed successfully")
	return nil
}

func (c *Tools) manualEmsdk() error {
	err := c.Shell.CheckCommand("git")

	if err != nil {
		return errors.New("git is not installed, please install git to proceed with emscripten installation")
	}

	err = c.Shell.GitClone("https://github.com/emscripten-core/emsdk.git", c.CloneDir)

	if err != nil {
		return fmt.Errorf("failed to clone emscripten repository: %w", err)
	}

	// join path with emsdk clone directory

	c.CloneDir = filepath.Join(c.CloneDir, "emsdk")

	err = c.Shell.RunScriptFileInDir("emsdk", c.CloneDir, "install", "latest")

	if err != nil {
		return fmt.Errorf("failed to run emscripten install script: %w", err)
	}

	err = c.Shell.RunScriptFileInDir("emsdk", c.CloneDir, "activate", "latest")

	if err != nil {
		return fmt.Errorf("failed to run emscripten activate script: %w", err)
	}

	return nil
}
