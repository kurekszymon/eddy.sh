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
	fmt.Println("Installing cmake using brew...")
	err := c.Shell.Brew("cmake")
	if err != nil {
		return err
	}
	fmt.Println("cmake installed successfully using brew")
	return nil
}

func (c *Tools) manualCmake() error {
	fmt.Println("Installing cmake using curl...")

	var url string
	if runtime.GOOS == "windows" {
		url = fmt.Sprintf("https://github.com/Kitware/CMake/releases/download/v%s/cmake-%s.zip", c.Cmake.Version, c.Cmake.Version)

	} else {
		url = fmt.Sprintf("https://github.com/Kitware/CMake/releases/download/v%s/cmake-%s.tar.gz", c.Cmake.Version, c.Cmake.Version)
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

	fmt.Println("SUCCESS: CMake installed successfully")
	return nil
}
