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

	fmt.Print("Do you want to proceed with this configuration: (Y/N) ")
	var i string
	fmt.Scan(&i)

	if i != "Y" && i != "y" {
		os.Exit(errors.WRONG_CONFIG)
	}
	fmt.Println("Proceeding with the installation...")

	if config.Platform.Brew {
		fmt.Println("Checking for brew...")
		err = handler.CheckCommand("brew")
		if err != nil {
			fmt.Println("Brew is not installed. Installing brew...")
			err = config.Installers.Tools.Brew.Install()
			if err != nil {
				fmt.Printf("ERROR: Failed to install brew: %v", err)
				fmt.Printf("Please try to install brew manually or specify manual installation in config.: %v", err)
				os.Exit(errors.BREW_SPECIFIED_BUT_NOT_INSTALLED)

			}
		}
		fmt.Println("Brew is installed and will be used for installation.")
	}

	err = handler.CheckCommand("git")
	if err != nil {
		fmt.Println("Git is not installed. Installing git...")
		err = config.Installers.Tools.Git.Install()
		if err != nil {
			fmt.Printf("ERROR: Failed to install git: %v", err)
		}
	}

	cpp := config.Installers.Cpp
	err = cpp.Cmake.Install()

	if err != nil {
		fmt.Printf("ERROR: Failed to install cmake: %v", err)
	}

	err = cpp.Ninja.Install()
	if err != nil {
		fmt.Printf("ERROR: Failed to install ninja: %v", err)
	}

}
