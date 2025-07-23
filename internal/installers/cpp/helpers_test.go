package cpp

import (
	"testing"

	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

func MockCppInstaller(t *testing.T) (*Installer, string) {
	t.Helper()

	tempDir := t.TempDir()

	mockShell := &installers.MockEddyDirShell{
		Shell:   shell.NewShellHandler(),
		TempDir: tempDir,
	}

	installer := &Installer{
		Shell:      mockShell,
		PkgManager: types.Manual,
		Available:  make(map[string]*installers.Tool),
	}

	return installer, tempDir
}
