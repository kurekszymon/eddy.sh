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

	ninjaUrl := fmt.Sprintf("https://github.com/ninja-build/ninja/releases/download/v%s/%s ", version, globals.NINJA_DIRNAME)
	eddyDir, err := c.Shell.GetEddyDir()
	if err != nil {
		return err
	}

	err = c.Shell.Curl(ninjaUrl, eddyDir)
	if err != nil {
		return fmt.Errorf("failed to download ninja: %w", err)
	}

	ninjaZipPath := filepath.Join(eddyDir, globals.NINJA_DIRNAME)
	err = c.Shell.Unzip(ninjaZipPath, eddyDir)
	if err != nil {
		return err
	}

	eddyBinDir, err := c.Shell.GetEddyBinDir()
	if err != nil {
		return err
	}

	ninjaPath := filepath.Join(eddyDir, "ninja")
	ninjaBinPath := filepath.Join(eddyBinDir, "ninja")
	os.Chmod(ninjaPath, 0755)

	c.Shell.Symlink(ninjaPath, ninjaBinPath)

	logger.Info("Ninja installed successfully")
	return nil
}
