package cpp

import (
	"errors"
	"path/filepath"
	"runtime"

	"github.com/kurekszymon/eddy.sh/internal/types"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

func (c *Tools) EmscriptenInstall() error {
	if c.PkgManager == types.Brew {
		return c.brewEmsdk()
	}

	return c.manualEmsdk()
}

func (c *Tools) brewEmsdk() error {
	utils.Log("Installing emscripten using brew", types.LogInfo)
	err := c.Shell.Brew("emscripten")
	if err != nil {
		return err
	}
	utils.Log("Emscripten installed successfully", types.LogInfo)
	return nil
}

func (c *Tools) manualEmsdk() error {
	err := c.Shell.CheckCommand("git")

	if err != nil {
		return errors.New(utils.FormatLogType("git is not installed. Please install git before proceeding with emscripten installation", types.LogError))
	}

	eddy_dir, err := c.Shell.GetEddyDir()
	if err != nil {
		return err
	}

	utils.Log("Cloning emscripten repository", types.LogInfo)
	err = c.Shell.GitClone("https://github.com/emscripten-core/emsdk.git", eddy_dir)

	if err != nil {
		return err
	}

	emsdk_dir := filepath.Join(eddy_dir, "emsdk")

	utils.Log("Running `emscripten install latest`", types.LogInfo)
	err = c.Shell.RunScriptFileInDir("emsdk", emsdk_dir, "install", "latest")

	if err != nil {
		return err
	}

	utils.Log("Running `emscripten activate latest`", types.LogInfo)
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

	utils.Log("Emscripten installed successfully", types.LogInfo)

	return nil
}
