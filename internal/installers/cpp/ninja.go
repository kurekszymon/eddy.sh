package cpp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kurekszymon/eddy.sh/internal/globals"
	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/types"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

func (c *Installer) NinjaInstall() error {
	if c.PkgManager == types.Brew {
		return c.brewNinja()
	}

	return c.manualNinja()
}

func (c *Installer) brewNinja() error {
	logger.Info("Installing ninja using brew")
	err := c.Shell.Brew("ninja")
	if err != nil {
		return err
	}
	logger.Info("Ninja installed successfully")
	return nil
}

// TODO: replace strings with typed constants
func (c *Installer) manualNinja() error {
	version, err := utils.DetermineVersion(c.Available["ninja"].Version, types.GHRepo{Name: "ninja", Owner: "ninja-build"})
	if err != nil {
		return err
	}

	logger.Info("Downloading ninja version " + version)

	ninja_url := fmt.Sprintf("https://github.com/ninja-build/ninja/releases/download/v%s/%s ", version, globals.NINJA_DIRNAME)
	err = c.Shell.Curl(ninja_url)
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
