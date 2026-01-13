package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const ConfigFileName = ".stamp.yaml"

type Config struct {
	Directory string `yaml:"directory"`
}

func DefaultConfig() *Config {
	return &Config{
		Directory: "docs/adr",
	}
}

func Load() (*Config, error) {
	configPath, err := FindConfigFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Save(dir string) error {
	configPath := filepath.Join(dir, ConfigFileName)

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func FindConfigFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		configPath := filepath.Join(dir, ConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", errors.New("no .stamp.yaml found (run 'stamp init' first)")
}

func (c *Config) ADRDirectory() (string, error) {
	configPath, err := FindConfigFile()
	if err != nil {
		return "", err
	}

	projectRoot := filepath.Dir(configPath)
	return filepath.Join(projectRoot, c.Directory), nil
}
