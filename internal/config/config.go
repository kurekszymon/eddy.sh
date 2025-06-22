package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/installers/cpp"
	"github.com/kurekszymon/eddy.sh/internal/installers/javascript"
	"github.com/kurekszymon/eddy.sh/internal/shell"
	"github.com/kurekszymon/eddy.sh/internal/types"
	"gopkg.in/yaml.v3"
)

type Installers struct {
	Cpp        *cpp.Tools
	Javascript *javascript.Tools
	Tools      *installers.Tools
}

type YamlPlatform struct {
	Brew   bool `yaml:"brew"`
	Manual bool `yaml:"manual_installation"`
}

type Config struct {
	Languages     []map[string][]map[string]string `yaml:"languages"`
	Git           types.Git                        `yaml:"git"`
	CustomScripts []map[string]string              `yaml:"custom_scripts"`
	Platform      YamlPlatform                     `yaml:"platform"`

	PkgManager types.PkgManager
	Installers *Installers
	Scripts    []types.CustomScript
}

// Process transforms the raw YAML structure into a more accessible format
func (c *Config) Process(shell *shell.ShellHandler) {
	platform := c.DetermineInstalationType(c.Platform)
	c.PkgManager = platform

	c.Installers = &Installers{
		Cpp:        &cpp.Tools{Shell: shell, PkgManager: c.PkgManager, CloneDir: c.Git.CloneDir},
		Javascript: &javascript.Tools{Shell: shell},
		Tools:      &installers.Tools{Shell: shell},
	}

	c.Installers.Tools = installers.GetTools(shell)

	var setter installers.ToolSetter

	for _, langGroup := range c.Languages {
		for langName, toolsList := range langGroup {
			for _, toolMap := range toolsList {
				for toolName, version := range toolMap {
					tool := &installers.Tool{
						Name:    toolName,
						Version: version,
					}

					switch strings.ToLower(langName) {
					case "cpp":
						setter = c.Installers.Cpp
					case "javascript":
						setter = c.Installers.Javascript
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

	if c.Installers.Cpp.Emscripten != nil {
		fmt.Printf("  Emscripten: %s (version: %s)\n", c.Installers.Cpp.Emscripten.Name, c.Installers.Cpp.Emscripten.Version)
	}
	if c.Installers.Cpp.Ninja != nil {
		fmt.Printf("  Ninja: %s (version: %s)\n", c.Installers.Cpp.Ninja.Name, c.Installers.Cpp.Ninja.Version)
	}
	if c.Installers.Cpp.Cmake != nil {
		fmt.Printf("  CMake: %s (version: %s)\n", c.Installers.Cpp.Cmake.Name, c.Installers.Cpp.Cmake.Version)
	}

	fmt.Printf("\nJavaScript Configuration:\n")
	if c.Installers.Javascript.Nvm != nil {
		fmt.Printf("  NVM: %s (version: %s)\n", c.Installers.Javascript.Nvm.Name, c.Installers.Javascript.Nvm.Version)
	}

	if c.Installers.Cpp.NotLoaded != nil && c.Installers.Javascript.NotLoaded != nil {
		fmt.Println("\nTools that won't be installed (please consider adding custom instructions for them):")
		for _, tool := range *c.Installers.Cpp.NotLoaded {
			fmt.Printf("- %s\n", tool.Name)
		}
		for _, tool := range *c.Installers.Javascript.NotLoaded {
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

	fmt.Printf("\nPlatform configuration:\n")
	fmt.Printf(" - use brew: %t\n", c.Platform.Brew)
	fmt.Printf(" - manual installation: %t\n", c.Platform.Manual)
	fmt.Println()
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

func (c *Config) DetermineInstalationType(platform YamlPlatform) types.PkgManager {
	switch {
	case platform.Brew:
		return types.Brew
	case platform.Manual:
		return types.Manual
	default:
		return types.Manual
	}
}
