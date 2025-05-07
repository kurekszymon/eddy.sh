package config

import (
	"os"
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/languages"
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

	// Processed fields for easier access
	LanguagesWrapper *Languages // Change this line
	Scripts          []types.CustomScript
}

// Process transforms the raw YAML structure into a more accessible format
func (c *Config) Process() {
	c.LanguagesWrapper = &Languages{
		Cpp:        &languages.CppTools{},
		Javascript: &languages.JsTools{},
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

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	config.Process()
	return config, nil
}
