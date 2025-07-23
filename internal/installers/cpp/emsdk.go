package cpp

import (
	"errors"
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

	eddy_dir, err := c.Shell.GetEddyDir()
	if err != nil {
		return err
	}

	logger.Info("Cloning emscripten repository")
	err = c.Shell.GitClone("https://github.com/emscripten-core/emsdk.git", eddy_dir)

	if err != nil {
		return err
	}

	emsdk_dir := filepath.Join(eddy_dir, "emsdk")

	logger.Info("Running `emscripten install latest`")
	err = c.Shell.RunScriptFileInDir("emsdk", emsdk_dir, "install", "latest")

	if err != nil {
		return err
	}

	logger.Info("Running `emscripten activate latest`")
	err = c.Shell.RunScriptFileInDir("emsdk", emsdk_dir, "activate", "latest")

	if err != nil {
		return err
	}

	logger.Info("Emscripten installed successfully")
	return nil
}
