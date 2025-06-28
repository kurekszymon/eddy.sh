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

	cmd := strings.ToLower(os.Args[1])

	switch cmd {
	case "install":
		if len(os.Args) < 3 {
			logger.Error("Please specify a tool to install, e.g. 'eddy.sh install nvm'")
			logger.Info("To install whole set of tools under language section run 'eddy.sh install javascript'")
			os.Exit(exit_codes.CLI_INSTALL_TOOL_NOT_SPECIFIED)
		}
		tool := strings.ToLower(os.Args[2])
		switch tool {

		// js
		case "javascript", "js":
			install(cfg.Installers.Javascript.Nvm, "nvm")
		case "nvm":
			install(cfg.Installers.Javascript.Nvm, "nvm")

		// c++
		case "cpp", "c++":
			install(cfg.Installers.Cpp.Cmake, "cmake")
			install(cfg.Installers.Cpp.Emscripten, "emscripten")
			install(cfg.Installers.Cpp.Ninja, "ninja")
		case "cmake":
			install(cfg.Installers.Cpp.Cmake, "cmake")
		case "emscripten":
			install(cfg.Installers.Cpp.Emscripten, "emscripten")
		case "ninja":
			install(cfg.Installers.Cpp.Ninja, "ninja")

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

func install(tool *installers.Tool, toolName string) {
	if tool != nil {
		err := tool.Install()
		if err != nil {
			msg := fmt.Sprintf("%s was not installed %s", tool.Name, err)
			logger.Error(msg)
			os.Exit(exit_codes.TOOL_NOT_INSTALLED)
		}
	} else {
		msg := fmt.Sprintf("%s is not present in your config file.", toolName)
		logger.Error(msg)
		os.Exit(exit_codes.CLI_INSTALL_TOOL_NOT_SPECIFIED)
	}
}

func printHelp() {
	fmt.Println(`
eddy.sh - Universal developer environment installer

Usage:
  eddy.sh install <tool>      Install a specific tool or tool group
  eddy.sh help                Show this help message

Examples:
  eddy.sh install nvm         Install Node Version Manager (nvm)
  eddy.sh install javascript  Install all JavaScript tools (e.g., nvm)
  eddy.sh install cmake       Install CMake
  eddy.sh install cpp         Install all C++ tools (cmake, emscripten, ninja)

Available tools:
  javascript, js     All JavaScript tools (currently: nvm)
  nvm                Node Version Manager
  cpp, c++           All C++ tools (cmake, emscripten, ninja)
  cmake              CMake build system
  emscripten         Emscripten compiler
  ninja              Ninja build system

For more information, see: https://github.com/kurekszymon/eddy.sh`)
}
