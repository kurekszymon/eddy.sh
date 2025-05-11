package main

import (
	"fmt"
	"os"

	"github.com/kurekszymon/eddy.sh/internal/config"
	"github.com/kurekszymon/eddy.sh/internal/errors"
	"github.com/kurekszymon/eddy.sh/internal/shell"
)

func main() {
	handler := &shell.ShellHandler{}

	config, err := config.LoadConfig("config.yaml", handler)
	if err != nil {
		// TODO: handle no config - generate sample
		fmt.Printf("Failed to load config: %v", err)
		os.Exit(errors.NO_CONFIG)
	}

	config.Print()
	var i string

	fmt.Print("Do you want to proceed with this configuration: (Y/N) ")
	fmt.Scan(&i)

	if i != "Y" && i != "y" {
		os.Exit(errors.WRONG_CONFIG)
	}
	fmt.Println("Proceeding with the installation...")
	// TODO: check for git.
	if config.Platform.Brew {
		fmt.Println("Checking for brew...")
		// install brew if not installed
		fmt.Println("Brew is installed and will be used for installation.")
	}

	cpp := config.LanguagesWrapper.Cpp
	err = cpp.Cmake.Install()

	if err != nil {
		fmt.Printf("ERROR: Failed to install cmake: %v", err)
	}

	err = cpp.Ninja.Install()
	if err != nil {
		fmt.Printf("ERROR: Failed to install ninja: %v", err)
	}
	// cpp.Emscripten.Install()

}
