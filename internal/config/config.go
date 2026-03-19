package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type Profile struct {
	Name       string `yaml:"name"`
	Email      string `yaml:"email"`
	SigningKey string `yaml:"signingkey,omitempty"`
	GPGSign    bool   `yaml:"gpgsign,omitempty"`
	SSHKey     string `yaml:"sshkey,omitempty"`
}

type Config struct {
	Profiles map[string]Profile `yaml:"profiles"`
}

func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".git-profiles.yaml")
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Profiles: make(map[string]Profile)}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}

	return &cfg, nil
}

func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}
