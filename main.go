package main

import (
	"os"

	"github.com/kurekszymon/eddy.sh/internal/config"
	"github.com/kurekszymon/eddy.sh/internal/error_codes"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/types"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

func main() {
	handler := &shell.ShellHandler{}

	config, err := config.LoadConfig("config.yaml", handler)
	if err != nil {
		// TODO: handle no config - generate sample
		utils.Log("Failed to load config, please check the config file or generate a new one.", types.LogError)
		os.Exit(error_codes.NO_CONFIG)
	}

	config.Print()

	utils.PromptConfirm("Do you want to proceed with this configuration: (Y/N) ", "ERROR: Failed to load config (user aborted)", error_codes.WRONG_CONFIG)
	utils.Log("Proceeding with the installation...", types.LogInfo)

	if config.Platform.Brew {
		utils.Log("Checking for brew...", types.LogInfo)
		err = handler.CheckCommand("brew")
		if err != nil {
			utils.Log("Brew is not installed. Installing brew...", types.LogWarning)
			err = config.Installers.Tools.Brew.Install()
			if err != nil {
				utils.Log("Failed to install brew", types.LogError)
				utils.Log("Please try to install brew manually or specify manual installation in config.", types.LogWarning)
				os.Exit(error_codes.BREW_SPECIFIED_BUT_NOT_INSTALLED)
			}
		}
		utils.Log("Brew is installed and will be used for installation.", types.LogInfo)
	}

	utils.Log("Checking for git...", types.LogInfo)
	err = handler.CheckCommand("git")
	if err != nil {
		utils.Log("Git is not installed. Installing git...", types.LogWarning)
		err = config.Installers.Tools.Git.Install()
		if err != nil {
			utils.Log("Failed to install git", types.LogError)
			os.Exit(error_codes.NO_GIT)
		}
	}

	// REPOSITORIES
	repos := config.Git.Repos
	if len(repos) > 0 {
		for _, repo := range repos {
			utils.Log("Cloning repository: "+repo, types.LogInfo)
			err = handler.GitClone(repo, config.Git.CloneDir)
			if err != nil {
				utils.Log("Failed to clone repository: "+repo, types.LogError)
			}
		}
	}

	// SCRIPTS
	custom_scripts := config.Scripts
	if len(custom_scripts) > 0 {
		for _, script := range custom_scripts {
			utils.Log("Running custom script: "+script.Name, types.LogInfo)
			err = handler.RunCustomScript(script.Command)
			if err != nil {
				utils.Log("Failed to run custom script: "+script.Name, types.LogError)
			}
		}
	}

	cpp := config.Installers.Cpp
	cpp.Install() // install now returns a map of errors, utilize this.

	js := config.Installers.Javascript
	js.Install()

	utils.Log("Please remember to add ~/.eddy.sh/bin to your PATH to access tools installed in the process.", types.LogWarning)
}
