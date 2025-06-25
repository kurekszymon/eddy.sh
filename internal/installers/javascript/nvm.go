package javascript

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/kurekszymon/eddy.sh/internal/logger"
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
	logger.Warn("Note: NVM is not available via Homebrew. It will be installed manually.")
	return c.manualNvm()
}

func (c *Tools) manualNvm() error {
	logger.Info("Downloading NVM version " + c.Nvm.Version)

	nvm_url := fmt.Sprintf("https://raw.githubusercontent.com/nvm-sh/nvm/v%s/install.sh", c.Nvm.Version)
	err := c.Shell.Curl(nvm_url)
	if err != nil {
		return err
	}

	eddy_dir, err := c.Shell.GetEddyDir()
	if err != nil {
		return err
	}

	err = c.Shell.RunScriptFileInDir("install.sh", eddy_dir)
	if err != nil {
		return err
	}

	return nil
}
