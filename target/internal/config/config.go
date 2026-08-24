package config

import (
	"os"
	"strconv"
)

// Config carries the runtime options for the BMS console server.
type Config struct {
	Addr     string
	DataDir  string
	SeedDemo bool
}

// Load reads configuration from environment variables with defaults.
func Load() Config {
	return Config{
		Addr:     envString("BMS_ADDR", ":18080"),
		DataDir:  envString("BMS_DATA", "./data"),
		SeedDemo: envBool("BMS_SEED", true),
	}
}

func envString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
