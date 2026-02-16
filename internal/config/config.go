package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

//go:embed default.yaml
var defaultConfig []byte

type Config struct {
	Storage    StorageConfig    `yaml:"storage"`
	Scheduler  SchedulerConfig  `yaml:"scheduler"`
	Copilot    CopilotConfig    `yaml:"copilot"`
	Image      ImageConfig      `yaml:"image"`
	Categories []CategoryConfig `yaml:"categories"`
}

type StorageConfig struct {
	Path string `yaml:"path"`
}

type SchedulerConfig struct {
	IntervalMinutes int `yaml:"interval_minutes"`
}

type CopilotConfig struct {
	Model string `yaml:"model"`
}

type ImageConfig struct {
	MaxWidth   int    `yaml:"max_width"`
	MaxFiles   int    `yaml:"max_files"`
	SaveImages bool   `yaml:"save_images"`
	Format     string `yaml:"format"`
}

type CategoryConfig struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Examples    []string `yaml:"examples"`
	Color       string   `yaml:"color"`
}

func Load(path string) (*Config, error) {
	resolvedPath, err := ResolvePath(path)
	if err != nil {
		return nil, err
	}
	path = resolvedPath

	// Initialize config if not exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := initConfig(path); err != nil {
			return nil, fmt.Errorf("init config: %w", err)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ResolvePath(path string) (string, error) {
	// Expand ~ to home directory
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get home dir: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	return filepath.Clean(path), nil
}

func Default() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(defaultConfig, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal default config: %w", err)
	}
	return &cfg, nil
}

func Write(path string, cfg *Config) error {
	resolvedPath, err := ResolvePath(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	dir := filepath.Dir(resolvedPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(resolvedPath, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

func initConfig(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	if err := os.WriteFile(path, defaultConfig, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}
