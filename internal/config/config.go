package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/languages"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/types"
	"gopkg.in/yaml.v3"
)

type Languages struct {
	Cpp        *languages.CppTools
	Javascript *languages.JsTools
}

type Config struct {
	Languages     []map[string][]map[string]string `yaml:"languages"`
	Git           types.Git                        `yaml:"git"`
	CustomScripts []map[string]string              `yaml:"custom_scripts"`
	Platform      struct {
		Brew bool `yaml:"brew"`
		Apt  bool `yaml:"apt"`
	} `yaml:"platform"`
	LanguagesWrapper *Languages
	Scripts          []types.CustomScript
}

// Process transforms the raw YAML structure into a more accessible format
func (c *Config) Process(shell *shell.ShellHandler) {
	c.LanguagesWrapper = &Languages{
		Cpp:        &languages.CppTools{Shell: shell},
		Javascript: &languages.JsTools{Shell: shell},
	}

	for _, langGroup := range c.Languages {
		for langName, toolsList := range langGroup {
			for _, toolMap := range toolsList {
				for toolName, version := range toolMap {
					tool := &languages.Tool{
						Name:    toolName,
						Version: version,
					}

					var setter languages.ToolSetter

					switch strings.ToLower(langName) {
					case "cpp":
						setter = c.LanguagesWrapper.Cpp
					case "javascript":
						setter = c.LanguagesWrapper.Javascript
					}

					if setter != nil {
						setter.SetTool(toolName, tool)
					}
				}
			}
		}
	}

	c.Scripts = make([]types.CustomScript, 0, len(c.CustomScripts))
	for _, scriptMap := range c.CustomScripts {
		for name, cmd := range scriptMap {
			c.Scripts = append(c.Scripts, types.CustomScript{Name: name, Command: cmd})
		}
	}
}

func (c *Config) Print() {
	fmt.Printf("C++ Configuration:\n")

	if c.LanguagesWrapper.Cpp.Emscripten != nil {
		fmt.Printf("  Emscripten: %s (version: %s)\n", c.LanguagesWrapper.Cpp.Emscripten.Name, c.LanguagesWrapper.Cpp.Emscripten.Version)
	}
	if c.LanguagesWrapper.Cpp.Ninja != nil {
		fmt.Printf("  Ninja: %s (version: %s)\n", c.LanguagesWrapper.Cpp.Ninja.Name, c.LanguagesWrapper.Cpp.Ninja.Version)
	}
	if c.LanguagesWrapper.Cpp.Cmake != nil {
		fmt.Printf("  CMake: %s (version: %s)\n", c.LanguagesWrapper.Cpp.Cmake.Name, c.LanguagesWrapper.Cpp.Cmake.Version)
	}

	fmt.Printf("\nJavaScript Configuration:\n")
	if c.LanguagesWrapper.Javascript.Nvm != nil {
		fmt.Printf("  NVM: %s (version: %s)\n", c.LanguagesWrapper.Javascript.Nvm.Name, c.LanguagesWrapper.Javascript.Nvm.Version)
	}

	if c.LanguagesWrapper.Cpp.NotLoaded != nil && c.LanguagesWrapper.Javascript.NotLoaded != nil {
		fmt.Println("\nTools that won't be installed (please consider adding custom instructions for them):")
		for _, tool := range *c.LanguagesWrapper.Cpp.NotLoaded {
			fmt.Printf("- %s\n", tool.Name)
		}
		for _, tool := range *c.LanguagesWrapper.Javascript.NotLoaded {
			fmt.Printf("- %s\n", tool.Name)
		}
	}

	fmt.Printf("\nGit Configuration:\n")
	fmt.Printf("  Clone Directory: %s\n", c.Git.CloneDir)
	fmt.Printf("  Repositories:\n")
	for _, repo := range c.Git.Repos {
		fmt.Printf("    - %s\n", repo)
	}

	fmt.Printf("\nCustom Scripts:\n")
	for _, script := range c.Scripts {
		fmt.Printf("  %s: %s\n", script.Name, script.Command)
	}
}

func LoadConfig(filename string, shell *shell.ShellHandler) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	config.Process(shell)
	return config, nil
}
