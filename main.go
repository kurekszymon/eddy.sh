package main

import (
	"fmt"
	"os"

	"github.com/kurekszymon/eddy.sh/internal/config"
	"github.com/kurekszymon/eddy.sh/internal/error_codes"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

func main() {
	handler := &shell.ShellHandler{}

	config, err := config.LoadConfig("config.yaml", handler)
	if err != nil {
		// TODO: handle no config - generate sample
		fmt.Printf("-- Failed to load config: %v", err)
		os.Exit(error_codes.NO_CONFIG)
	}

	config.Print()

	utils.PromptConfirm("Do you want to proceed with this configuration: (Y/N) ", "ERROR: Failed to load config (user aborted)", error_codes.WRONG_CONFIG)
	fmt.Println("-- Proceeding with the installation...")

	if config.Platform.Brew {
		fmt.Println("-- Checking for brew...")
		err = handler.CheckCommand("brew")
		if err != nil {
			fmt.Println("-- Brew is not installed. Installing brew...")
			err = config.Installers.Tools.Brew.Install()
			if err != nil {
				fmt.Printf("ERROR: Failed to install brew: %v", err)
				fmt.Printf("Please try to install brew manually or specify manual installation in config.: %v", err)
				os.Exit(error_codes.BREW_SPECIFIED_BUT_NOT_INSTALLED)
			}
		}
		fmt.Println("-- Brew is installed and will be used for installation.")
	}

	err = handler.CheckCommand("git")
	if err != nil {
		fmt.Println("-- Git is not installed. Installing git...")
		err = config.Installers.Tools.Git.Install()
		if err != nil {
			fmt.Printf("-- ERROR: Failed to install git: %v", err)
			os.Exit(error_codes.NO_GIT)
		}
	}

	// REPOSITORIES
	repos := config.Git.Repos
	if len(repos) > 0 {
		for _, repo := range repos {
			fmt.Printf("Cloning repository: %s\n", repo)
			err = handler.GitClone(repo, config.Git.CloneDir)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// SCRIPTS
	custom_scripts := config.Scripts
	if len(custom_scripts) > 0 {
		for _, script := range custom_scripts {
			fmt.Println("-- Running custom script:", script.Name)
			err = handler.RunCustomScript(script.Command)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	cpp := config.Installers.Cpp
	cpp.Install() // maybe should have "print errors" method to print all errors at the end
}
