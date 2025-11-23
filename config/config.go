package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration
type Config struct {
	API          APIConfig          `yaml:"api"`
	Server       ServerConfig       `yaml:"server"`
	EuropeanPads EuropeanPadsConfig `yaml:"european_pads"`
}

// APIConfig holds API-related settings
type APIConfig struct {
	BaseURL            string `yaml:"base_url"`
	Limit              int    `yaml:"limit"`
	Mode               string `yaml:"mode"`
	Format             string `yaml:"format"`
	Key                string `yaml:"key"`                  // Stargazer API key from environment
	UseProductionImage bool   `yaml:"use_production_image"` // Replace -dev with -prod in image URLs
}

// ServerConfig holds server settings
type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

// EuropeanPadsConfig holds European launch pad filtering settings
type EuropeanPadsConfig struct {
	Enabled      bool     `yaml:"enabled"`
	LocationIDs  []int    `yaml:"location_ids"`
	CountryCodes []string `yaml:"country_codes"`
}

var (
	instance *Config
	once     sync.Once
)

// Get returns the singleton config instance
func Get() *Config {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}

// loadConfig loads configuration from YAML file
func loadConfig() *Config {
	cfg := &Config{
		// Default values
		API: APIConfig{
			BaseURL: "https://lldev.thespacedevs.com/2.3.0/launches/upcoming/",
			Limit:   5,
			Mode:    "list",
			Format:  "json",
			Key:     "", // Will be loaded from env
		},
		Server: ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
		EuropeanPads: EuropeanPadsConfig{
			Enabled: false,
		},
	}

	// Try to load from config.yaml
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		// Config file not found, use defaults
		return cfg
	}

	// Parse YAML
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		fmt.Printf("Warning: failed to parse config.yaml: %v\n", err)
		return cfg
	}

	// Load RLEU_API_KEY from environment (overrides YAML config)
	if apiKey := os.Getenv("RLEU_API_KEY"); apiKey != "" {
		cfg.API.Key = apiKey
	}

	return cfg
}

// GetAPIURL builds the full API URL with query parameters
func (c *Config) GetAPIURL() string {
	return fmt.Sprintf("%s?limit=%d&mode=%s&format=%s",
		c.API.BaseURL,
		c.API.Limit,
		c.API.Mode,
		c.API.Format,
	)
}

// GetServerAddress returns the server address as host:port
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}
