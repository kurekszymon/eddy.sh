package main

import (
	"fmt"
	"log"

	"github.com/kurekszymon/eddy.sh/internal/config"
)

func main() {
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		// TODO: handle no config - generate sample
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("C++ Configuration:\n")

	if config.LanguagesWrapper.Cpp.Emscripten != nil {
		fmt.Printf("  Emscripten: %s (version: %s)\n", config.LanguagesWrapper.Cpp.Emscripten.Name, config.LanguagesWrapper.Cpp.Emscripten.Version)
		errored := config.LanguagesWrapper.Cpp.Emscripten.Install()
		if errored != nil {
			fmt.Println(errored.Error())
		}
	}
	if config.LanguagesWrapper.Cpp.Ninja != nil {
		fmt.Printf("  Ninja: %s (version: %s)\n", config.LanguagesWrapper.Cpp.Ninja.Name, config.LanguagesWrapper.Cpp.Ninja.Version)
	}
	if config.LanguagesWrapper.Cpp.Cmake != nil {
		fmt.Printf("  CMake: %s (version: %s)\n", config.LanguagesWrapper.Cpp.Cmake.Name, config.LanguagesWrapper.Cpp.Cmake.Version)
	}

	fmt.Printf("\nJavaScript Configuration:\n")
	if config.LanguagesWrapper.Javascript.Nvm != nil {
		fmt.Printf("  NVM: %s (version: %s)\n", config.LanguagesWrapper.Javascript.Nvm.Name, config.LanguagesWrapper.Javascript.Nvm.Version)
	}

	fmt.Printf("\nGit Configuration:\n")
	fmt.Printf("  Clone Directory: %s\n", config.Git.CloneDir)
	fmt.Printf("  Repositories:\n")
	for _, repo := range config.Git.Repos {
		fmt.Printf("    - %s\n", repo)
	}

	fmt.Printf("\nCustom Scripts:\n")
	for _, script := range config.Scripts {
		fmt.Printf("  %s: %s\n", script.Name, script.Command)
	}

}
