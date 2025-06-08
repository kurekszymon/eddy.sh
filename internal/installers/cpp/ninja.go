package cpp

import (
	"fmt"
	"path/filepath"

	"github.com/kurekszymon/eddy.sh/internal/globals"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

func (c *Tools) NinjaInstall() error {
	if c.PkgManager == types.Brew {
		return c.brewNinja()
	}

	return c.manualNinja()
}

func (c *Tools) brewNinja() error {
	fmt.Println("Installing ninja...")
	err := c.Shell.Brew("ninja")
	if err != nil {
		return fmt.Errorf("failed to install ninja: %w", err)
	}
	fmt.Println("ninja installed successfully")
	return nil
}

func (c *Tools) manualNinja() error {
	fmt.Println("Installing ninja manually...")
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
		return fmt.Errorf("failed to get eddy dir: %w", err)
	}

	ninja := filepath.Join(eddy_dir, "ninja")
	c.Shell.ChmodX(ninja)

	return nil
}
