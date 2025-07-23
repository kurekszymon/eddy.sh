package main

import (
	"maps"
	"os"
	"path"

	"github.com/kurekszymon/eddy.sh/internal/cli"
	"github.com/kurekszymon/eddy.sh/internal/config"
	"github.com/kurekszymon/eddy.sh/internal/exit_codes"
	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/installers/cpp"
	"github.com/kurekszymon/eddy.sh/internal/installers/general"
	"github.com/kurekszymon/eddy.sh/internal/installers/javascript"
	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

type Installers struct {
	Cpp        *cpp.Installer
	Javascript *javascript.Installer
	Tools      *general.Installer
}

func main() {
	handler := shell.NewShellHandler()

	cli.HandleArgs(handler)

	configFile := determineConfigFile(handler)

	yaml, err := config.Load(shell.ExpandPath(configFile))

	if err != nil {
		logger.Error("Failed to load config, please check " + configFile)
		os.Exit(exit_codes.WRONG_CONFIG)
	}

	cfg := config.Build(yaml)

	specificInstallers := map[string]installers.Installer{
		"cpp":        &cpp.Installer{Shell: handler, PkgManager: cfg.PkgManager},
		"javascript": &javascript.Installer{Shell: handler, PkgManager: cfg.PkgManager},
	}

	generalInstaller := general.NewGeneralInstaller(handler, cfg.PkgManager)

	config.PrintConfig(cfg)
	utils.PromptConfirm("Do you want to proceed with this configuration?", "ERROR: Failed to load config (user aborted)", exit_codes.WRONG_CONFIG)
	logger.Info("Proceeding with the installation...")

	logger.Info("Using config file: " + configFile)
	generalInstallerErrors := generalInstaller.Install()

	if generalInstallerErrors["brew"] != nil {
		logger.Error("Failed to install brew")
		logger.Warn("Please try to install brew manually or specify manual installation in config.")
		os.Exit(exit_codes.BREW_SPECIFIED_BUT_NOT_INSTALLED)
	}

	if generalInstallerErrors["git"] != nil {
		if err != nil {
			logger.Error("Failed to install git")
			os.Exit(exit_codes.NO_GIT)
		}
	}

	// REPOSITORIES
	logger.Info("Preparing to run clone repositories: ")
	for _, repo := range cfg.Git.Repos {
		logger.Info(" - " + repo)
	}

	logger.Warn("If you plan to use SSH authentication with GitHub, GitLab, or Bitbucket, make sure to generate SSH key and add it to your account:")
	logger.Info("GitHub:    https://docs.github.com/en/authentication/connecting-to-github-with-ssh")
	logger.Info("GitLab:    https://docs.gitlab.com/user/ssh/")
	logger.Info("Bitbucket: https://support.atlassian.com/bitbucket-cloud/docs/set-up-an-ssh-key/")
	utils.PromptConfirm("Please continue only after you make sure you've added SSH keys to your account - otherwise 'git clone' may fail.", "Git installation denied by the user.", exit_codes.SSH_KEYS_DENIED)

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
	for k, v := range cfg.Tools {
		for _, tool := range v {
			specificInstallers[k].SetTool(tool.Name, &tool)
		}
	}

	errors := make(map[string]map[string]error)

	for name, installer := range specificInstallers {
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

func determineConfigFile(handler *shell.ShellHandler) string {
	var configFile string
	flags := utils.HandleFlags()

	if flags[utils.Config] != "" {
		configFile = flags[utils.Config]
	} else {
		eddy_dir, err := handler.GetEddyDir()

		if err != nil {
			logger.Error("Failed to get Eddy directory: " + err.Error())
			os.Exit(exit_codes.SOMETHING_WENT_WRONG)
		}

		configFile = path.Join(eddy_dir, "config.yaml")

		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			logger.Warn("Config file not found at " + configFile)
			utils.PromptConfirm("Do you want to use the default config? [Y/n]", "User denied using default config.", exit_codes.NO_CONFIG)
			handler.Curl("https://raw.githubusercontent.com/kurekszymon/eddy.sh/refs/heads/main/config.yaml")

			determineConfigFile(handler)
		}
	}

	return configFile
}
