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
			installFunc: func() error { return shell.Brew("git") },
		},
		Brew: &Tool{
			Name:        "brew",
			Version:     "latest",
			installFunc: func() error { return shell.Curl("https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh") },
		},
	}
}
