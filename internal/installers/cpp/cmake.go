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
	var cmake_dir string
	var cmake_bin_path string
	var url string
	version, err := utils.DetermineVersion(c.Available["cmake"].Version, types.GHRepo{Name: "CMake", Owner: "Kitware"})

	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		cmake_dir = fmt.Sprintf("cmake-%s-windows-arm64", version)
		cmake_bin_path = filepath.Join(cmake_dir, "bin")

		url = fmt.Sprintf("https://github.com/Kitware/CMake/releases/download/v%s/%s.zip", version, cmake_dir)
	} else {
		cmake_dir = fmt.Sprintf("cmake-%s-macos-universal", version)
		cmake_bin_path = filepath.Join(cmake_dir, "Cmake.app", "Contents", "bin")
		url = fmt.Sprintf("https://github.com/Kitware/CMake/releases/download/v%s/%s.tar.gz", version, cmake_dir)
	}

	err = c.Shell.Curl(url)
	if err != nil {
		return err
	}

	filename := filepath.Base(url)
	err = c.Shell.Unzip(filename, "")
	if err != nil {
		return err
	}

	eddy_dir, err := c.Shell.GetEddyDir()
	if err != nil {
		return err
	}

	cmake_bin := filepath.Join(eddy_dir, cmake_bin_path)

	for _, bin := range []string{"cmake", "cpack", "ctest", "ccmake"} {
		c.Shell.Symlink(filepath.Join(cmake_bin, bin), bin)
	}

	logger.Info("CMake installed successfully")
	return nil
}
