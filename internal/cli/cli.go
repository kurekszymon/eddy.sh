package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/exit_codes"
	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/installers/cpp"
	"github.com/kurekszymon/eddy.sh/internal/installers/javascript"
	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

func HandleArgs(handler *shell.ShellHandler) {
	if len(os.Args) < 2 {
		// No command, run interactive installer (main logic)
		return
	}

	cmd := strings.ToLower(os.Args[1])

	switch cmd {
	case "install":
		if len(os.Args) < 3 {
			logger.Error("Please specify a tool to install, e.g. 'eddy.sh install nvm'")
			logger.Info("To install whole set of tools under language section run 'eddy.sh install javascript'")
			os.Exit(exit_codes.CLI_INSTALL_TOOL_NOT_SPECIFIED)
		}
		tool := strings.ToLower(os.Args[2])
		version := "latest"
		if len(os.Args) > 3 {
			version = strings.ToLower(os.Args[3])
		}

		jsInstaller := &javascript.Installer{Shell: handler, PkgManager: types.Manual}
		cppInstaller := &cpp.Installer{Shell: handler, PkgManager: types.Manual}

		switch tool {
		// js
		case "javascript", "js":
			loadTool(jsInstaller, "nvm", "latest")
			install(jsInstaller, "javascript")
		case "nvm":
			loadTool(jsInstaller, "nvm", version)
			install(jsInstaller, "nvm")

		// c++
		case "cpp", "c++":
			loadTool(cppInstaller, "cmake", "latest")
			loadTool(cppInstaller, "emscripten", "latest")
			loadTool(cppInstaller, "ninja", "latest")
			install(cppInstaller, "c++")
		case "cmake":
			loadTool(cppInstaller, "cmake", version)
			install(cppInstaller, "cmake")
		case "emscripten":
			loadTool(cppInstaller, "emscripten", version)
			install(cppInstaller, "emscripten")
		case "ninja":
			loadTool(cppInstaller, "ninja", version)
			install(cppInstaller, "ninja")

		default:
			logger.Error("Unknown tool " + tool)
			os.Exit(exit_codes.UNKNOWN_TOOL)
		}

	case "help":
		printHelp()
		os.Exit(exit_codes.SUCCESS)

	default:
		fmt.Println("Unknown command: " + cmd)
		os.Exit(exit_codes.UNKNOWN_COMMAND)
	}

	os.Exit(exit_codes.SUCCESS)
}

func install(installer installers.Installer, name string) {
	errors := installer.Install()
	if len(errors) > 0 {
		for toolName, err := range errors {
			msg := fmt.Sprintf("%s was not installed: %s", toolName, err)
			logger.Error(msg)
		}
		os.Exit(exit_codes.TOOL_NOT_INSTALLED)
	}
}

func loadTool(group installers.Installer, name string, version string) {
	tool := &installers.Tool{
		Name:    name,
		Version: version,
	}
	group.SetTool(name, tool)
}

func printHelp() {
	fmt.Println(`
eddy.sh - Universal developer environment installer

Usage:
  eddy.sh install <tool> [version]   Install a specific tool or tool group (optionally specify version)
  eddy.sh help                       Show this help message

Examples:
  eddy.sh install nvm                Install Node Version Manager (nvm) (latest version)
  eddy.sh install nvm 0.40.3         Install Node Version Manager (nvm) version 0.40.3
  eddy.sh install javascript         Install all JavaScript tools (e.g., nvm)
  eddy.sh install cmake              Install CMake (latest version)
  eddy.sh install cmake 3.27.0       Install CMake version 3.27.0
  eddy.sh install cpp                Install all C++ tools (cmake, emscripten, ninja)

Available tools:
  javascript, js     All JavaScript tools (currently: nvm)
  nvm                Node Version Manager
  cpp, c++           All C++ tools (cmake, emscripten, ninja)
  cmake              CMake build system
  emscripten         Emscripten compiler
  ninja              Ninja build system

You can specify a version for any tool. If omitted, the latest version will be installed.

For more information, see: https://github.com/kurekszymon/eddy.sh`)
}
