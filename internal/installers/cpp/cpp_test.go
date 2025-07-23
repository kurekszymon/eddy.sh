package cpp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kurekszymon/eddy.sh/internal/installers"
)

func Test_CmakeManualInstall(t *testing.T) {
	t.Parallel()

	installer, tempDir := MockCppInstaller(t)
	tool := &installers.Tool{Name: "cmake", Version: "latest"}
	installer.SetTool("cmake", tool)

	errors := installer.Install()

	if len(errors) > 0 {
		t.Fatalf("Cmake.Install() error = %v", errors)
	}

	binaries := []string{"cmake", "ccmake", "ctest"}
	for _, bin := range binaries {
		expectedFile := filepath.Join(tempDir, "bin", bin)
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist, but it doesn't", expectedFile)
		}
	}
}

func Test_EmsdkManualInstall(t *testing.T) {
	t.Parallel()

	installer, tempDir := MockCppInstaller(t)
	tool := &installers.Tool{Name: "emscripten", Version: "latest"}
	installer.SetTool("emscripten", tool)

	errors := installer.Install()

	if len(errors) > 0 {
		t.Fatalf("Emscripten.Install() error = %v", errors)
	}

	binaries := []string{"emsdk", "emsdk_env.sh"}

	for _, bin := range binaries {
		expectedFile := filepath.Join(tempDir, "emsdk", bin)
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist, but it doesn't", expectedFile)
		}
	}
}

func Test_NinjaManualInstall(t *testing.T) {
	t.Parallel()

	installer, tempDir := MockCppInstaller(t)
	tool := &installers.Tool{Name: "ninja", Version: "latest"}
	installer.SetTool("ninja", tool)

	errors := installer.Install()

	if len(errors) > 0 {
		t.Fatalf("Cmake.Install() error = %v", errors)
	}

	expectedFile := filepath.Join(tempDir, "bin", "ninja")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it doesn't", expectedFile)
	}
}
