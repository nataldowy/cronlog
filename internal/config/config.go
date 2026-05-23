package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level cronlog configuration.
type Config struct {
	Store   StoreConfig   `yaml:"store"`
	Notify  NotifyConfig  `yaml:"notify"`
	Runner  RunnerConfig  `yaml:"runner"`
}

// StoreConfig configures the log entry storage backend.
type StoreConfig struct {
	Path string `yaml:"path"`
}

// NotifyConfig configures failure notification settings.
type NotifyConfig struct {
	SMTPHost   string   `yaml:"smtp_host"`
	SMTPPort   int      `yaml:"smtp_port"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
}

// RunnerConfig configures default runner behaviour.
type RunnerConfig struct {
	DefaultTimeout time.Duration `yaml:"default_timeout"`
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: invalid: %w", err)
	}

	return &cfg, nil
}

// validate checks required fields and applies defaults.
func (c *Config) validate() error {
	if c.Store.Path == "" {
		c.Store.Path = "/var/log/cronlog"
	}
	if c.Notify.SMTPPort == 0 {
		c.Notify.SMTPPort = 25
	}
	if c.Runner.DefaultTimeout == 0 {
		c.Runner.DefaultTimeout = 24 * time.Hour
	}
	return nil
}
