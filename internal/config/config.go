package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/installers"
	"github.com/kurekszymon/eddy.sh/internal/types"
	"gopkg.in/yaml.v3"
)

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

type Config struct {
	PkgManager    types.PkgManager
	Tools         map[string][]installers.Tool
	Git           types.Git
	CustomScripts []types.CustomScript
}

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

func Build(rawConfig *RawConfig) *Config {
	pkgManager := determinePackageManager(rawConfig.Platform)
	tools := processTools(rawConfig.Languages)
	scripts := processScripts(rawConfig.CustomScripts)

	return &Config{
		PkgManager:    pkgManager,
		Tools:         tools,
		Git:           types.Git{CloneDir: rawConfig.Git.CloneDir, Repos: rawConfig.Git.Repos},
		CustomScripts: scripts,
	}
}

func Print(cfg *Config) {
	fmt.Println("--- Configuration ---")
	fmt.Printf("Package Manager: %s\n", cfg.PkgManager)

	if len(cfg.Git.Repos) > 0 {
		fmt.Println("\nGit:")
		fmt.Printf("  Clone Directory: %s\n", cfg.Git.CloneDir)
		fmt.Println("  Repositories:")
		for _, repo := range cfg.Git.Repos {
			fmt.Printf("    - %s\n", repo)
		}
	}

	if len(cfg.Tools) > 0 {
		fmt.Println("\nTools:")
		for lang, tools := range cfg.Tools {
			fmt.Printf("  %s:\n", lang)
			for _, tool := range tools {
				fmt.Printf("    - %s (version: %s)\n", tool.Name, tool.Version)
			}
		}
	}

	if len(cfg.CustomScripts) > 0 {
		fmt.Println("\nCustom Scripts:")
		for _, script := range cfg.CustomScripts {
			fmt.Printf("  - Name: %s, Command: %s\n", script.Name, script.Command)
		}
	}
	fmt.Println("---------------------")
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
