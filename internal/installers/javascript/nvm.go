package javascript

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/types"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

func (j *Installer) NvmInstall() error {
	if runtime.GOOS == "windows" {
		return errors.New("NVM installation is not yet supported on Windows. Please use NVM for Windows")
	}

	if j.PkgManager == types.Brew {
		return j.brewNvm()
	}

	return j.manualNvm()
}

func (j *Installer) brewNvm() error {
	logger.Warn("Note: NVM is not available via Homebrew. It will be installed manually.")
	return j.manualNvm()
}

func (j *Installer) manualNvm() error {
	version, err := utils.DetermineVersion(j.Available["nvm"].Version, types.GHRepo{Name: "nvm", Owner: "nvm-sh"})
	if err != nil {
		return err
	}

	logger.Info("Downloading NVM version " + j.Available["nvm"].Version)

	eddy_dir, err := j.Shell.GetEddyDir()
	if err != nil {
		return err
	}
	nvm_url := fmt.Sprintf("https://raw.githubusercontent.com/nvm-sh/nvm/v%s/install.sh", version)
	err = j.Shell.Curl(nvm_url, eddy_dir)
	if err != nil {
		return err
	}

	err = j.Shell.RunScriptFileInDir("install.sh", eddy_dir)
	if err != nil {
		return err
	}

	return nil
}
