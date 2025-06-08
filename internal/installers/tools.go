package installers

import (
	"github.com/kurekszymon/eddy.sh/internal/shell"
)

type Tools struct {
	Shell *shell.ShellHandler

	Git  *Tool
	Brew *Tool
}

func GetTools(shell *shell.ShellHandler) *Tools {
	return &Tools{
		Shell: shell,
		Git: &Tool{
			Name:        "git",
			Version:     "latest",
			InstallFunc: func() error { return shell.Brew("git") },
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
