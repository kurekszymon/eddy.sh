package installers

import (
	"path/filepath"

	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/utils"
)

type MockEddyDirShell struct {
	shell.Shell
	TempDir string
}

func (m *MockEddyDirShell) GetEddyDir() (string, error) {
	return m.TempDir, nil
}

func (m *MockEddyDirShell) GetEddyBinDir() (string, error) {
	dir := filepath.Join(m.TempDir, "bin")
	utils.EnsureDir(dir)
	return dir, nil
}
