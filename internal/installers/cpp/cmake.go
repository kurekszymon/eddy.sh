package cpp

import (
	"fmt"

	"github.com/kurekszymon/eddy.sh/internal/types"
)

func (c *Tools) CmakeInstall() error {
	if c.PkgManager == types.Brew {
		return c.brewCmake()
	}

	return c.manualCmake()
}

func (c *Tools) brewCmake() error {
	fmt.Println("Installing cmake using brew...")
	err := c.Shell.Brew("cmake")
	if err != nil {
		return fmt.Errorf("failed to install cmake using brew: %w", err)
	}
	fmt.Println("cmake installed successfully using brew")
	return nil
}

func (c *Tools) manualCmake() error {
	fmt.Println("Manual installation of cmake is not supported yet.")
	fmt.Println("Please follow the instructions at https://cmake.org/install/")
	return nil
}
