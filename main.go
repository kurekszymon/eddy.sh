package main

import (
	"os"

	"github.com/kurekszymon/eddy.sh/internal/cli"
	"github.com/kurekszymon/eddy.sh/internal/config"
	"github.com/kurekszymon/eddy.sh/internal/exit_codes"
	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

func main() {
	handler := &shell.ShellHandler{}

	config := config.Init(handler)
	cli.HandleArgs(handler, config)

	err := config.Load(handler)
	if err != nil {
		logger.Error("Failed to load config, please check " + config.File)
		os.Exit(exit_codes.WRONG_CONFIG)
	}

	config.Print()

	utils.PromptConfirm("Do you want to proceed with this configuration?", "ERROR: Failed to load config (user aborted)", exit_codes.WRONG_CONFIG)
	logger.Info("Proceeding with the installation...")

	if config.Platform.Brew {
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

	err = handler.CheckCommand("git")
	if err != nil {
		logger.Warn("Git is not installed. Installing git...")
		err = config.Installers.Tools.Git.Install()
		if err != nil {
			logger.Error("Failed to install git")
			os.Exit(exit_codes.NO_GIT)
		}
	}

	logger.Warn("If you plan to use SSH with GitHub, GitLab, or Bitbucket, make sure to generate SSH key and add it to your account:")
	logger.Info("GitHub:    https://docs.github.com/en/authentication/connecting-to-github-with-ssh")
	logger.Info("GitLab:    https://docs.gitlab.com/user/ssh/")
	logger.Info("Bitbucket: https://support.atlassian.com/bitbucket-cloud/docs/set-up-an-ssh-key/")
	utils.PromptConfirm("Please continue only after you make sure you've added SSH keys to your account - otherwise 'git clone' may fail.", "Git installation denied by the user.", exit_codes.SSH_KEYS_DENIED)

	// REPOSITORIES
	repos := config.Git.Repos
	if len(repos) > 0 {
		for _, repo := range repos {
			logger.Info("Cloning repository: " + repo)
			err = handler.GitClone(repo, config.Git.CloneDir)
			if err != nil {
				logger.Error(err.Error())
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
				logger.Error(err.Error())
			}
		}
	}

	// INSTALLERS
	cpp := config.Installers.Cpp
	cpp_errors := cpp.Install()

	js := config.Installers.Javascript
	js_errors := js.Install()

	utils.PrintInstallErrors(cpp_errors, js_errors)

	logger.Warn("Please remember to add ~/.eddy.sh/bin to your PATH to access tools installed in the process.")
}
