package javascript

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/kurekszymon/eddy.sh/internal/types"
	"github.com/kurekszymon/eddy.sh/internal/utils"
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
	utils.Log("Note: NVM is not available via Homebrew. It will be installed manually.", types.LogWarning)
	return c.manualNvm()
}

func (c *Tools) manualNvm() error {
	utils.Log("Downloading NVM version "+c.Nvm.Version, types.LogInfo)

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
