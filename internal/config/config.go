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
	Cpp        *cpp.Installer
	Javascript *javascript.Installer
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
	File       string
}

// RawConfig is a struct that directly maps to the structure of the config.yaml file.
// It is used only for parsing and should not be used directly by the application.	ype RawConfig struct {
type RawConfig struct {
	Languages []map[string][]map[string]string `yaml:"languages"`
	Git       struct {
		CloneDir string   `yaml:"clone_dir"`
		Repos    []string `yaml:"repos"`
	} `yaml:"git"`
	CustomScripts []map[string]string `yaml:"custom_scripts"`
	Platform      struct {
		Brew   bool `yaml:"brew"`
		Manual bool `yaml:"manual_installation"`
	} `yaml:"platform"`
}

// Load parses the configuration file at the given path and returns a RawConfig struct.
// It is the responsibility of the caller to handle file not found errors and to
// provide a default configuration if necessary.
func Load(path string) (*RawConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var rawConfig RawConfig
	err = yaml.Unmarshal(data, &rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &rawConfig, nil
}

type Settings struct {
	PkgManager    types.PkgManager
	Tools         map[string][]installers.Tool
	Git           types.Git
	CustomScripts []types.CustomScript
}

// Build creates a new Settings object from a RawConfig.
func Build(rawConfig *RawConfig) *Settings {
	pkgManager := determinePackageManager(rawConfig.Platform)
	tools := processTools(rawConfig.Languages)
	scripts := processScripts(rawConfig.CustomScripts)

	return &Settings{
		PkgManager:    pkgManager,
		Tools:         tools,
		Git:           types.Git{CloneDir: rawConfig.Git.CloneDir, Repos: rawConfig.Git.Repos},
		CustomScripts: scripts,
	}
}

func determinePackageManager(platform struct {
	Brew   bool `yaml:"brew"`
	Manual bool `yaml:"manual_installation"`
}) types.PkgManager {
	if platform.Brew {
		return types.Brew
	}
	return types.Manual
}

func processTools(languages []map[string][]map[string]string) map[string][]installers.Tool {
	toolsByLanguage := make(map[string][]installers.Tool)
	for _, langGroup := range languages {
		for langName, toolsList := range langGroup {
			lowerLangName := strings.ToLower(langName)
			if _, ok := toolsByLanguage[lowerLangName]; !ok {
				toolsByLanguage[lowerLangName] = []installers.Tool{}
			}
			for _, toolMap := range toolsList {
				for toolName, version := range toolMap {
					tool := installers.Tool{
						Name:    toolName,
						Version: version,
					}
					toolsByLanguage[lowerLangName] = append(toolsByLanguage[lowerLangName], tool)
				}
			}
		}
	}
	return toolsByLanguage
}

func processScripts(customScripts []map[string]string) []types.CustomScript {
	var scripts []types.CustomScript
	for _, scriptMap := range customScripts {
		for name, cmd := range scriptMap {
			scripts = append(scripts, types.CustomScript{Name: name, Command: cmd})
		}
	}
	return scripts
}

func Init(shell *shell.ShellHandler) *Config {
	config := &Config{}

	platform := config.DetermineInstalationType(config.Platform)
	config.PkgManager = platform

	config.Installers = &Installers{
		Cpp:        &cpp.Installer{Shell: shell, PkgManager: config.PkgManager, CloneDir: config.Git.CloneDir},
		Javascript: &javascript.Installer{Shell: shell, PkgManager: config.PkgManager},
		Tools:      &installers.Tools{Shell: shell},
	}

	// config.Installers.Tools = installers.GetTools(shell)

	return config
}

// func (c *Config) ProcessYaml(shell *shell.ShellHandler) {
// 	var setter installers.ToolSetter

// 	for _, langGroup := range c.Languages {
// 		for langName, toolsList := range langGroup {
// 			for _, toolMap := range toolsList {
// 				for toolName, version := range toolMap {
// 					tool := &installers.Tool{
// 						Name:    toolName,
// 						Version: version,
// 					}

// 					switch strings.ToLower(langName) {
// 					case "cpp":
// 						setter = c.Installers.Cpp
// 					case "javascript":
// 						setter = c.Installers.Javascript
// 					}

// 					if setter != nil {
// 						setter.SetTool(toolName, tool)
// 					}
// 				}
// 			}
// 		}
// 	}

// 	c.Scripts = make([]types.CustomScript, 0, len(c.CustomScripts))
// 	for _, scriptMap := range c.CustomScripts {
// 		for name, cmd := range scriptMap {
// 			c.Scripts = append(c.Scripts, types.CustomScript{Name: name, Command: cmd})
// 		}
// 	}
// }

// func (c *Config) Print() {
// 	// do it nicely
// 	fmt.Printf("C++ Configuration:\n")

// 	if c.Installers.Cpp.Emscripten != nil {
// 		fmt.Printf("  Emscripten: %s (version: %s)\n", c.Installers.Cpp.Emscripten.Name, c.Installers.Cpp.Emscripten.Version)
// 	}
// 	if c.Installers.Cpp.Ninja != nil {
// 		fmt.Printf("  Ninja: %s (version: %s)\n", c.Installers.Cpp.Ninja.Name, c.Installers.Cpp.Ninja.Version)
// 	}
// 	if c.Installers.Cpp.Cmake != nil {
// 		fmt.Printf("  CMake: %s (version: %s)\n", c.Installers.Cpp.Cmake.Name, c.Installers.Cpp.Cmake.Version)
// 	}

// 	fmt.Printf("\nJavaScript Configuration:\n")
// 	if c.Installers.Javascript.Nvm != nil {
// 		fmt.Printf("  NVM: %s (version: %s)\n", c.Installers.Javascript.Nvm.Name, c.Installers.Javascript.Nvm.Version)
// 	}

// 	if c.Installers.Cpp.NotLoaded != nil && c.Installers.Javascript.NotLoaded != nil {
// 		fmt.Println("\nTools that won't be installed (please consider adding custom instructions for them):")
// 		for _, tool := range *c.Installers.Cpp.NotLoaded {
// 			fmt.Printf("- %s\n", tool.Name)
// 		}
// 		for _, tool := range *c.Installers.Javascript.NotLoaded {
// 			fmt.Printf("- %s\n", tool.Name)
// 		}
// 	}

// 	fmt.Printf("\nGit Configuration:\n")
// 	fmt.Printf("  Clone Directory: %s\n", c.Git.CloneDir)
// 	fmt.Printf("  Repositories:\n")
// 	for _, repo := range c.Git.Repos {
// 		fmt.Printf("    - %s\n", repo)
// 	}

// 	fmt.Printf("\nCustom Scripts:\n")
// 	for _, script := range c.Scripts {
// 		fmt.Printf("  %s: %s\n", script.Name, script.Command)
// 	}

// 	fmt.Printf("\nPlatform configuration:\n")
// 	fmt.Printf(" - use brew: %t\n", c.Platform.Brew)
// 	fmt.Printf(" - manual installation: %t\n", c.Platform.Manual)
// 	fmt.Println()
// }

// func (config *Config) Load(shell *shell.ShellHandler) error {
// 	config_path := config.File

// 	if config_path == "" {
// 		eddy_dir, err := shell.GetEddyDir()
// 		if err != nil {
// 			logger.Error("Something went horribly wrong. Please report an issue.") // paste link to issues
// 			os.Exit(exit_codes.SOMETHING_WENT_WRONG)
// 		}
// 		config.File = globals.CONFIG_FILE
// 		config_path = path.Join(eddy_dir, globals.CONFIG_FILE)
// 	}

// 	_, err := os.Stat(config_path)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			logger.Warn(config_path + " was not found.")
// 			prompt := fmt.Sprintf("Do you want to download default config from %s?", globals.CONFIG_URL)
// 			utils.PromptConfirm(prompt, "User denied downloading config.", exit_codes.NO_CONFIG)
// 			shell.Curl("https://raw.githubusercontent.com/kurekszymon/eddy.sh/refs/heads/main/config.yaml")
// 		}
// 	}

// 	data, err := os.ReadFile(config_path)
// 	if err != nil {
// 		return err
// 	}

// 	err = yaml.Unmarshal(data, config)
// 	if err != nil {
// 		return err
// 	}

// 	config.ProcessYaml()

// 	return nil
// }

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
