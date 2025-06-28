package main

import (
	"os"

	"github.com/kurekszymon/eddy.sh/internal/cli"
	"github.com/kurekszymon/eddy.sh/internal/config"
	"github.com/kurekszymon/eddy.sh/internal/exit_codes"
	"github.com/kurekszymon/eddy.sh/internal/globals"
	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

func main() {
	handler := &shell.ShellHandler{}

	config, err := config.LoadConfig(handler)
	if err != nil {
		logger.Error("Failed to load config, please check ~/.eddy.sh" + globals.CONFIG_FILE)
		os.Exit(exit_codes.WRONG_CONFIG)
	}

	cli.HandleArgs(handler, config)

	config.Print()

	utils.PromptConfirm("Do you want to proceed with this configuration?", "ERROR: Failed to load config (user aborted)", exit_codes.WRONG_CONFIG)
	logger.Info("Proceeding with the installation...")

	if config.Platform.Brew {
		logger.Info("Checking for brew...")
		err = handler.CheckCommand("brew")
		if err != nil {
			logger.Warn("Brew is not installed. Installing brew...")
			err = config.Installers.Tools.Brew.Install()
			if err != nil {
				logger.Error("Failed to install brew")
				logger.Warn("Please try to install brew manually or specify manual installation in config.")
				os.Exit(exit_codes.BREW_SPECIFIED_BUT_NOT_INSTALLED)
			}
		}
		logger.Info("Brew is installed and will be used for installation.")
	}

	logger.Info("Checking for git...")
	err = handler.CheckCommand("git")
	if err != nil {
		logger.Warn("Git is not installed. Installing git...")
		err = config.Installers.Tools.Git.Install()
		if err != nil {
			logger.Error("Failed to install git")
			os.Exit(exit_codes.NO_GIT)
		}
	}

	// REPOSITORIES
	repos := config.Git.Repos
	if len(repos) > 0 {
		for _, repo := range repos {
			logger.Info("Cloning repository: " + repo)
			err = handler.GitClone(repo, config.Git.CloneDir)
			if err != nil {
				logger.Error("Failed to clone repository: " + repo)
			}
		}
	}

	// SCRIPTS
	custom_scripts := config.Scripts
	if len(custom_scripts) > 0 {
		for _, script := range custom_scripts {
			logger.Info("Running custom script: " + script.Name)
			err = handler.RunCustomScript(script.Command)
			if err != nil {
				logger.Error("Failed to run custom script: " + script.Name)
			}
		}
	}

	cpp := config.Installers.Cpp
	cpp.Install() // install now returns a map of errors, utilize this.

	js := config.Installers.Javascript
	js.Install()

	logger.Warn("Please remember to add ~/.eddy.sh/bin to your PATH to access tools installed in the process.")
}
