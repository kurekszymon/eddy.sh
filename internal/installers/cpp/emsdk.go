package cpp

import (
	"errors"
	"path/filepath"
	"runtime"

	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

func (c *Tools) EmscriptenInstall() error {
	if c.PkgManager == types.Brew {
		return c.brewEmsdk()
	}

	return c.manualEmsdk()
}

func (c *Tools) brewEmsdk() error {
	logger.Info("Installing emscripten using brew")
	err := c.Shell.Brew("emscripten")
	if err != nil {
		return err
	}
	logger.Info("Emscripten installed successfully")
	return nil
}

func (c *Tools) manualEmsdk() error {
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

	if runtime.GOOS == "windows" {
		c.Shell.Symlink(emsdk_dir, "emsdk.sh")
		c.Shell.Symlink(emsdk_dir, "emsdk.bat")
		c.Shell.Symlink(emsdk_dir, "emsdk.ps1")
		c.Shell.Symlink(emsdk_dir, "emsdk_env.ps1")
		c.Shell.Symlink(emsdk_dir, "emsdk_env.bat")
		c.Shell.Symlink(emsdk_dir, "emsdk_env.sh")

	} else {
		c.Shell.Symlink(emsdk_dir, "emsdk")
		c.Shell.Symlink(emsdk_dir, "emsdk.sh")
		c.Shell.Symlink(emsdk_dir, "emsdk_env.sh")
	}

	logger.Info("Emscripten installed successfully")

	return nil
}
