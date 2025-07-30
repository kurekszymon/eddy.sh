package cpp

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

func (c *Installer) EmscriptenInstall() error {
	if c.PkgManager == types.Brew {
		return c.brewEmsdk()
	}

	return c.manualEmsdk()
}

func (c *Installer) brewEmsdk() error {
	logger.Info("Installing emscripten using brew")
	err := c.Shell.Brew("emscripten")
	if err != nil {
		return err
	}
	logger.Info("Emscripten installed successfully")
	return nil
}

func (c *Installer) manualEmsdk() error {
	err := c.Shell.CheckCommand("git")
	if err != nil {
		return errors.New("git is not installed. Please install git before proceeding with emscripten installation")
	}
	version := c.Available["emscripten"].Version

	eddyDir, err := c.Shell.GetEddyDir()
	if err != nil {
		return err
	}

	logger.Info("Cloning emscripten repository")
	err = c.Shell.GitClone("https://github.com/emscripten-core/emsdk.git", eddyDir)

	if err != nil {
		return err
	}

	emsdkDir := filepath.Join(eddyDir, "emsdk")

	msg := fmt.Sprintf("Running `emscripten install %s`", version)
	logger.Info(msg)
	err = c.Shell.RunScriptFileInDir("emsdk", emsdkDir, "install", version)

	if err != nil {
		return err
	}

	msg = fmt.Sprintf("Running `emscripten activate %s`", version)
	logger.Info(msg)
	err = c.Shell.RunScriptFileInDir("emsdk", emsdkDir, "activate", version)

	if err != nil {
		return err
	}

	logger.Info("Emscripten installed successfully")
	return nil
}
