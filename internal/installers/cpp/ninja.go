package cpp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kurekszymon/eddy.sh/internal/globals"
	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

func (c *Tools) NinjaInstall() error {
	if c.PkgManager == types.Brew {
		return c.brewNinja()
	}

	return c.manualNinja()
}

func (c *Tools) brewNinja() error {
	logger.Info("Installing ninja using brew")
	err := c.Shell.Brew("ninja")
	if err != nil {
		return err
	}
	logger.Info("Ninja installed successfully")
	return nil
}

func (c *Tools) manualNinja() error {
	logger.Info("Downloading ninja version " + c.Ninja.Version)

	ninja_url := fmt.Sprintf("https://github.com/ninja-build/ninja/releases/download/v%s/%s ", c.Ninja.Version, globals.NINJA_DIRNAME)
	err := c.Shell.Curl(ninja_url)
	if err != nil {
		return fmt.Errorf("failed to download ninja: %w", err)
	}

	err = c.Shell.Unzip(globals.NINJA_DIRNAME, "")
	if err != nil {
		return err
	}

	eddy_dir, err := c.Shell.GetEddyDir()
	if err != nil {
		return err
	}

	ninja_path := filepath.Join(eddy_dir, "ninja")
	os.Chmod(ninja_path, 0755)

	c.Shell.Symlink(ninja_path, "ninja")

	logger.Info("Ninja installed successfully")
	return nil
}
