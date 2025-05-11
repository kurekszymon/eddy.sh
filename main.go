package main

import (
	"fmt"
	"log"

	"github.com/kurekszymon/eddy.sh/internal/config"
	"github.com/kurekszymon/eddy.sh/internal/shell"
)

func main() {
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		// TODO: handle no config - generate sample
		log.Fatalf("Failed to load config: %v", err)
	}

	config.Print()

	cpp := config.LanguagesWrapper.Cpp

	cpp.Emscripten.Install()

	handler := &shell.ShellHandler{Verbose: true}

	err = handler.Run("ls . && ls 1>&2 .")
	if err != nil {
		fmt.Println("Error:", err)
	}

}
