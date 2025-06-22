package cpp

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/kurekszymon/eddy.sh/internal/types"
)

func (c *Tools) CmakeInstall() error {
	if c.PkgManager == types.Brew {
		return c.brewCmake()
	}

	return c.manualCmake()
}

func (c *Tools) brewCmake() error {
	fmt.Println("-- Installing cmake using brew...")
	err := c.Shell.Brew("cmake")
	if err != nil {
		return err
	}
	fmt.Println("-- cmake installed successfully using brew")
	return nil
}

func (c *Tools) manualCmake() error {
	fmt.Println("-- Installing cmake using curl...")
	var cmake_dir string
	var cmake_bin_path string
	var url string
	if runtime.GOOS == "windows" {
		cmake_dir = fmt.Sprintf("cmake-%s-windows-arm64", c.Cmake.Version)
		cmake_bin_path = filepath.Join(cmake_dir, "bin")

		url = fmt.Sprintf("https://github.com/Kitware/CMake/releases/download/v%s/%s.zip", c.Cmake.Version, cmake_dir)
	} else {
		cmake_dir = fmt.Sprintf("cmake-%s-macos-universal", c.Cmake.Version)
		cmake_bin_path = filepath.Join(cmake_dir, "Cmake.app", "Contents", "bin")
		url = fmt.Sprintf("https://github.com/Kitware/CMake/releases/download/v%s/%s.tar.gz", c.Cmake.Version, cmake_dir)
	}

	err := c.Shell.Curl(url)
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

	fmt.Println("-- Creating symlinks for cmake...")

	cmake_bin := filepath.Join(eddy_dir, cmake_bin_path)

	c.Shell.Symlink(filepath.Join(cmake_bin, "cmake"), "cmake")
	c.Shell.Symlink(filepath.Join(cmake_bin, "cpack"), "cpack")
	c.Shell.Symlink(filepath.Join(cmake_bin, "ctest"), "ctest")
	c.Shell.Symlink(filepath.Join(cmake_bin, "ccmake"), "ccmake")

	fmt.Println("-- SUCCESS: CMake installed successfully")
	return nil
}
