package javascript

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/kurekszymon/eddy.sh/internal/types"
)

func (c *Tools) NvmInstall() error {
	if runtime.GOOS == "windows" {
		return errors.New("NVM installation is not yet supported on Windows. Please use NVM for Windows")
	}

	if c.PkgManager == types.Brew {
		return c.brewNvm()
	}

	return c.manualNvm()
}

func (c *Tools) brewNvm() error {
	fmt.Println("Note: NVM is not available via Homebrew. It will be installed manually.")
	return c.manualNvm()
}

func (c *Tools) manualNvm() error {
	fmt.Println("-- Installing NVM manually...")

	nvm_url := fmt.Sprintf("https://raw.githubusercontent.com/nvm-sh/nvm/v%s/install.sh", c.Nvm.Version)
	err := c.Shell.Curl(nvm_url)
	if err != nil {
		return err
	}

	eddy_dir, err := c.Shell.GetEddyDir()
	if err != nil {
		return fmt.Errorf("failed to get eddy directory: %w", err)
	}

	err = c.Shell.RunScriptFileInDir("install.sh", eddy_dir)
	if err != nil {
		return fmt.Errorf("failed to run NVM install script: %w", err)
	}

	return nil
}
