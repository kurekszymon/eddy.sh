package installers

import (
	"errors"
	"runtime"

	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

type GeneralTools struct {
	Shell *shell.ShellHandler

	Git  *Tool
	Brew *Tool
}

type Tools struct {
	Shell      *shell.ShellHandler
	PkgManager types.PkgManager
	CloneDir   string

	Available map[string]*Tool
	NotLoaded []Tool
}

func GetTools(shell *shell.ShellHandler) *GeneralTools {
	return &GeneralTools{
		Git: &Tool{
			Name:    "git",
			Version: "latest",
			InstallFunc: func() error {
				if runtime.GOOS == "windows" {
					// move to seperate file probably
					return errors.New("git installation for windows is not supported yet. Please follow manual steps from https://git-scm.com/downloads/win")
				}
				err := shell.Brew("git")
				if err != nil {
					return err
				}

				return nil
			},
		},
		Brew: &Tool{
			Name:    "brew",
			Version: "latest",
			InstallFunc: func() error {
				err := shell.Curl("https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh")
				if err != nil {
					return err
				}

				eddy_dir, err := shell.GetEddyDir()
				if err != nil {
					return err
				}

				err = shell.RunScriptFileInDir("install.sh", eddy_dir)

				if err != nil {
					return err
				}

				return nil
			},
		},
	}
}
