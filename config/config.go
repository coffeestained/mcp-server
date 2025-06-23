package config

import "github.com/spf13/viper"

type Config struct {
	Server        ServerConfig
	Github        GithubConfig
	StackExchange StackExchangeConfig
	OpenAPI       OpenAPIConfig
}

type ServerConfig struct {
	Port string
}

// GithubConfig is now optional. We will check if APIKey is present.
type GithubConfig struct {
	APIKey       string
	Repositories map[string]string
}

// StackExchangeConfig is also optional.
type StackExchangeConfig struct {
	APIKey string
}

// OpenAPIConfig now supports multiple named schemas.
type OpenAPIConfig struct {
	Schemas map[string]string // Maps a friendly name to a schema path (URL or file)
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// It's okay if the config file doesn't exist, we can run with defaults.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error and run with empty config
		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Set a default port if not provided
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	return &cfg, nil
}