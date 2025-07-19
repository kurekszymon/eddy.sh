package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/config"
	"github.com/kurekszymon/eddy.sh/internal/exit_codes"
	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/shell"
)

func HandleArgs(handler *shell.ShellHandler, cfg *config.Config) {
	if len(os.Args) < 2 {
		// No command, run interactive installer (main logic)
		return
	}

	flags_present := false
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "--config" || os.Args[i] == "-c" && i+1 < len(os.Args) {
			cfg.File = shell.ExpandPath(os.Args[i+1])
			logger.Info("Using config file: " + cfg.File)
			flags_present = true
			continue
		}
		// parse --platform / --pkgManager
	}
	if flags_present {
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

		switch tool {

		// js
		case "javascript", "js":
			loadTool(cfg.Installers.Javascript, "nvm", "latest")

			install(cfg.Installers.Javascript.Nvm)
		case "nvm":
			loadTool(cfg.Installers.Javascript, "nvm", version)

			install(cfg.Installers.Javascript.Nvm)

		// c++
		case "cpp", "c++":
			loadTool(cfg.Installers.Cpp, "cmake", "latest")
			loadTool(cfg.Installers.Cpp, "emscripten", "latest")
			loadTool(cfg.Installers.Cpp, "cmake", "latest")

			install(cfg.Installers.Cpp.Cmake)
			install(cfg.Installers.Cpp.Emscripten)
			install(cfg.Installers.Cpp.Ninja)
		case "cmake":
			loadTool(cfg.Installers.Cpp, "cmake", version)
			install(cfg.Installers.Cpp.Cmake)
		case "emscripten":
			loadTool(cfg.Installers.Cpp, "emscripten", version)
			install(cfg.Installers.Cpp.Emscripten)
		case "ninja":
			loadTool(cfg.Installers.Cpp, "ninja", version)

			install(cfg.Installers.Cpp.Ninja)

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

func install(tool *installers.Tool) {
	err := tool.Install()
	if err != nil {
		msg := fmt.Sprintf("%s was not installed %s", tool.Name, err)
		logger.Error(msg)
		os.Exit(exit_codes.TOOL_NOT_INSTALLED)
	}

}

func loadTool(group installers.ToolSetter, name string, version string) {
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


Flags (can be placed after the main command):
  --config, -c <file>        Use a custom config file (default: ~/.eddy.sh/config.yaml)

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
