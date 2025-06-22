package cpp

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/kurekszymon/eddy.sh/internal/types"
)

func (c *Tools) EmscriptenInstall() error {
	if c.PkgManager == types.Brew {
		return c.brewEmsdk()
	}

	return c.manualEmsdk()
}

func (c *Tools) brewEmsdk() error {
	fmt.Println("-- Installing emscripten with brew...")
	err := c.Shell.Brew("emscripten")
	if err != nil {
		return err
	}
	fmt.Println("-- emscripten installed successfully")
	return nil
}

func (c *Tools) manualEmsdk() error {
	err := c.Shell.CheckCommand("git")

	if err != nil {
		return errors.New("-- git is not installed, please install git to proceed with emscripten installation")
	}

	eddy_dir, err := c.Shell.GetEddyDir()
	if err != nil {
		return fmt.Errorf("-- failed to get eddy directory: %w", err)
	}

	fmt.Println("-- Cloning emscripten repository...")
	err = c.Shell.GitClone("https://github.com/emscripten-core/emsdk.git", eddy_dir)

	if err != nil {
		return err
	}

	emsdk_dir := filepath.Join(eddy_dir, "emsdk")

	fmt.Println("-- Running emscripten install script...")
	err = c.Shell.RunScriptFileInDir("emsdk", emsdk_dir, "install", "latest")

	if err != nil {
		return err
	}

	fmt.Println("-- Activating emscripten environment...")
	err = c.Shell.RunScriptFileInDir("emsdk", emsdk_dir, "activate", "latest")

	if err != nil {
		return err
	}

	fmt.Println("-- Creating symlinks for emscripten...")

	if runtime.GOOS == "windows" {
		c.Shell.Symlink(emsdk_dir, "emsdk.sh")
		c.Shell.Symlink(emsdk_dir, "emsdk.bat")
		c.Shell.Symlink(emsdk_dir, "emsdk.ps1")
		c.Shell.Symlink(emsdk_dir, "emsdk_env.ps1")
		c.Shell.Symlink(emsdk_dir, "emsdk_env.bat")
		c.Shell.Symlink(emsdk_dir, "emsdk_env.sh")

	} else {
		c.Shell.Symlink(emsdk_dir, "emsdk")
		c.Shell.Symlink(emsdk_dir, "emsdk_env.sh")
	}

	fmt.Println("-- SUCCESS: Emscripten installed successfully")

	return nil
}
