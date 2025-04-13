package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	MatchType  string `yaml:"match_type"`
	SubType    string `yaml:"sub_type"`
	NamePrefix string `yaml:"name_prefix"`
	NameRegex  string `yaml:"name_regex"`
	IgnoreType bool   `yaml:"ignore_type"`
	OutputFile string `yaml:"output_file"`
}

type Config struct {
	Rules []Rule `yaml:"rules"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

func DefaultConfig() *Config {
	return &Config{
		Rules: []Rule{},
	}
}

func (c *Config) MatchRule(matchType, subType, name string) string {
	for _, rule := range c.Rules {
		if !rule.IgnoreType && rule.MatchType != matchType {
			continue
		}

		if rule.SubType != "" && rule.SubType != subType {
			continue
		}

		if rule.NamePrefix != "" && (name == "" || !strings.HasPrefix(name, rule.NamePrefix)) {
			continue
		}

		if rule.NameRegex != "" {
			if name == "" {
				continue
			}
			matched, err := regexp.MatchString(rule.NameRegex, name)
			if err != nil || !matched {
				continue
			}
		}

		return rule.OutputFile
	}

	switch matchType {
	case "resource", "data":
		return fmt.Sprintf("%s_%s.tf", matchType, subType)
	case "module":
		return "modules.tf"
	case "variable":
		return "variables.tf"
	case "output":
		return "outputs.tf"
	case "locals":
		return "locals.tf"
	case "provider":
		return "providers.tf"
	default:
		return fmt.Sprintf("%s.tf", matchType)
	}
}
