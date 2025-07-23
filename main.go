package main

import (
	"maps"

	"github.com/kurekszymon/eddy.sh/internal/config"
	"github.com/kurekszymon/eddy.sh/internal/exit_codes"
	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/installers/cpp"
	"github.com/kurekszymon/eddy.sh/internal/installers/javascript"
	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

type Installers struct {
	Cpp        *cpp.Installer
	Javascript *javascript.Installer
	Tools      *installers.Tools
}

func main() {
	handler := shell.NewShellHandler()

	yaml, err := config.Load("config.yaml")

	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	cfg := config.Build(yaml)

	// Extract parseFlags from HandleArgs
	// use it to set a file.
	// fix a File to also be present in config.Settings
	// cli.HandleArgs(handler, cfg)

	// Handle file not exist -> download default config
	// err := config.Load(handler)
	// if err != nil {
	// 	logger.Error("Failed to load config, please check " + config.File)
	// 	os.Exit(exit_codes.WRONG_CONFIG)
	// }

	// Print a config
	// config.Print()
	// utils.PromptConfirm("Do you want to proceed with this configuration?", "ERROR: Failed to load config (user aborted)", exit_codes.WRONG_CONFIG)
	// logger.Info("Proceeding with the installation...")

	// handle brew and git installation
	// if config.Platform.Brew {
	// 	err = handler.CheckCommand("brew")
	// 	if err != nil {
	// 		logger.Warn("Brew is not installed. Installing brew...")
	// 		err = config.Installers.Tools.Brew.Install()
	// 		if err != nil {
	// 			logger.Error("Failed to install brew")
	// 			logger.Warn("Please try to install brew manually or specify manual installation in config.")
	// 			os.Exit(exit_codes.BREW_SPECIFIED_BUT_NOT_INSTALLED)
	// 		}
	// 	}
	// 	logger.Info("Brew is installed and will be used for installation.")
	// }

	// err = handler.CheckCommand("git")
	// if err != nil {
	// 	logger.Warn("Git is not installed. Installing git...")
	// 	err = config.Installers.Tools.Git.Install()
	// 	if err != nil {
	// 		logger.Error("Failed to install git")
	// 		os.Exit(exit_codes.NO_GIT)
	// 	}
	// }

	logger.Warn("If you plan to use SSH with GitHub, GitLab, or Bitbucket, make sure to generate SSH key and add it to your account:")
	logger.Info("GitHub:    https://docs.github.com/en/authentication/connecting-to-github-with-ssh")
	logger.Info("GitLab:    https://docs.gitlab.com/user/ssh/")
	logger.Info("Bitbucket: https://support.atlassian.com/bitbucket-cloud/docs/set-up-an-ssh-key/")
	utils.PromptConfirm("Please continue only after you make sure you've added SSH keys to your account - otherwise 'git clone' may fail.", "Git installation denied by the user.", exit_codes.SSH_KEYS_DENIED)

	// REPOSITORIES
	for _, repo := range cfg.Git.Repos {
		logger.Info("Cloning repository: " + repo)
		err = handler.GitClone(repo, cfg.Git.CloneDir)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	// SCRIPTS
	for _, script := range cfg.CustomScripts {
		logger.Info("Running custom script: " + script.Name)
		err = handler.RunCustomScript(script.Command)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	// LANGUAGES
	installers := map[string]installers.Installer{
		"cpp":        &cpp.Installer{Shell: handler, PkgManager: cfg.PkgManager, CloneDir: cfg.Git.CloneDir},
		"javascript": &javascript.Installer{Shell: handler, PkgManager: cfg.PkgManager, CloneDir: cfg.Git.CloneDir},
		// "tools":      &installers.Tools{Shell: handler, PkgManager: cfg.PkgManager, CloneDir: cfg.Git.CloneDir},
		// handle tools
	}

	for k, v := range cfg.Tools {
		for _, tool := range v {
			installers[k].SetTool(tool.Name, &tool)
		}
	}

	errors := make(map[string]map[string]error)

	for name, installer := range installers {
		installErrs := installer.Install()

		if len(installErrs) > 0 {
			if errors[name] == nil {
				errors[name] = make(map[string]error)
			}

			maps.Copy(errors[name], installErrs)
		}
	}

	utils.PrintInstallErrors(errors)

	logger.Warn("Please remember to add ~/.eddy.sh/bin to your PATH to access tools installed in the process.")
}
