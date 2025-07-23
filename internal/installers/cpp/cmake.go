package cpp

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/types"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

func (c *Installer) CmakeInstall() error {
	if c.PkgManager == types.Brew {
		return c.brewCmake()
	}

	return c.manualCmake()
}

func (c *Installer) brewCmake() error {
	logger.Info("Installing cmake using brew")
	err := c.Shell.Brew("cmake")
	if err != nil {
		return err
	}
	logger.Info("CMake installed successfully")
	return nil
}

func (c *Installer) manualCmake() error {
	logger.Info("Installing CMake manually")
	var cmakeDir string
	var cmakeBinPath string
	var url string
	version, err := utils.DetermineVersion(c.Available["cmake"].Version, types.GHRepo{Name: "CMake", Owner: "Kitware"})

	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		cmakeDir = fmt.Sprintf("cmake-%s-windows-arm64", version)
		cmakeBinPath = filepath.Join(cmakeDir, "bin")

		url = fmt.Sprintf("https://github.com/Kitware/CMake/releases/download/v%s/%s.zip", version, cmakeDir)
	} else {
		cmakeDir = fmt.Sprintf("cmake-%s-macos-universal", version)
		cmakeBinPath = filepath.Join(cmakeDir, "Cmake.app", "Contents", "bin")
		url = fmt.Sprintf("https://github.com/Kitware/CMake/releases/download/v%s/%s.tar.gz", version, cmakeDir)
	}

	eddyDir, err := c.Shell.GetEddyDir()
	if err != nil {
		return err
	}
	eddyBinDir, err := c.Shell.GetEddyBinDir()
	if err != nil {
		return err
	}

	err = c.Shell.Curl(url, eddyDir)
	if err != nil {
		return err
	}

	filename := filepath.Base(url)
	archivePath := filepath.Join(eddyDir, filename)
	err = c.Shell.Unzip(archivePath, eddyDir)
	if err != nil {
		return err
	}

	cmakeBin := filepath.Join(eddyDir, cmakeBinPath)

	for _, bin := range []string{"cmake", "cpack", "ctest", "ccmake"} {
		err := c.Shell.Symlink(filepath.Join(cmakeBin, bin), filepath.Join(eddyBinDir, bin))
		if err != nil {
			msg := fmt.Sprintf("Failed to create symlink for %s: %v", bin, err)
			logger.Error(msg)
			// grab tools that errored and return
		}
	}

	logger.Info("CMake installed successfully")
	return nil
}
