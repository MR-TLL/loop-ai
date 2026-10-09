package config

import (
	"fmt"
	"log"
	"os"
	_ "embed"

	"github.com/goccy/go-yaml"
)

//go:embed config.yaml
var defaultConfigBytes []byte

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Booking  BookingConfig  `yaml:"booking"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host                string `yaml:"host"`
	Port                int    `yaml:"port"`
	User                string `yaml:"user"`
	Password            string `yaml:"password"`
	Name                string `yaml:"name"`
	SSLMode             string `yaml:"sslmode"`
	MaxOpenConns        int    `yaml:"max_open_conns"`
	MaxIdleConns        int    `yaml:"max_idle_conns"`
	ConnMaxLifetimeMins int    `yaml:"conn_max_lifetime_mins"`
}

// DSN builds the postgres connection string from the DatabaseConfig fields.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

type BookingConfig struct {
	HoldExpirySecs         int `yaml:"hold_expiry_secs"`
	ExpiryPollIntervalSecs int `yaml:"expiry_poll_interval_secs"`
}

// Load reads config.yaml (embedded as default), then overlays an external
// config file if CONFIG_PATH env var is set, allowing per-environment overrides
// without rebuilding the binary.
func Load() *Config {
	cfg := &Config{}

	// Start from embedded defaults
	if err := yaml.Unmarshal(defaultConfigBytes, cfg); err != nil {
		log.Fatalf("config: failed to parse embedded config.yaml: %v", err)
	}

	// Optional external override via CONFIG_PATH env var
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("config: failed to read config file at %s: %v", path, err)
		}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			log.Fatalf("config: failed to parse config file at %s: %v", path, err)
		}
		log.Printf("config: loaded from %s", path)
	}

	return cfg
}
