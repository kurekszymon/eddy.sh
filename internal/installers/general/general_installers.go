package general

import (
	"errors"
	"runtime"
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

type Installer struct {
	Shell      shell.Shell
	PkgManager types.PkgManager

	Available map[string]*installers.Tool
	NotLoaded []installers.Tool
}

func NewGeneralInstaller(shell shell.Shell, packageManager types.PkgManager) *Installer {

	installer := &Installer{
		Shell:      shell,
		PkgManager: packageManager,
		Available:  make(map[string]*installers.Tool),
	}

	installer.SetTool("git", &installers.Tool{
		Name:    "git",
		Version: "latest",
	})

	installer.SetTool("brew", &installers.Tool{
		Name:    "brew",
		Version: "latest",
	})

	return installer
}

func (g *Installer) SetTool(toolName string, tool *installers.Tool) {
	toolKey := strings.ToLower(toolName)

	switch toolKey {
	case "git":

		tool.InstallFunc = func() error {
			err := g.Shell.CheckCommand("git")
			if err == nil {
				return nil
			}

			logger.Warn("Git is not installed. Installing git...")

			if runtime.GOOS == "windows" {
				return errors.New("git installation for windows is not supported yet. Please follow manual steps from https://git-scm.com/downloads/win")
			}

			return g.Shell.Brew("git")
		}

	case "brew":
		tool.InstallFunc = func() error {
			if g.PkgManager == "brew" {
				err := g.Shell.CheckCommand("brew")
				if err == nil {
					return nil
				}
				eddyDir, err := g.Shell.GetEddyDir()
				if err != nil {
					return err
				}
				logger.Warn("Brew is not installed. Installing brew...")
				err = g.Shell.Curl("https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh", eddyDir)
				if err != nil {
					return err
				}
				return g.Shell.RunScriptFileInDir("install.sh", eddyDir)
			}
			return nil
		}

	default:
		g.NotLoaded = append(g.NotLoaded, *tool)
		return
	}

	g.Available[toolKey] = tool
}

func (g *Installer) Install() map[string]error {
	errors := make(map[string]error)

	for name, tool := range g.Available {
		if err := tool.Install(); err != nil {
			errors[name] = err

			logger.Error("Failed to install " + name + ": " + err.Error())
		}
	}

	return errors
}
