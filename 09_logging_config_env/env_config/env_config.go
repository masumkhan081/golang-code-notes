// env_config.go
// Demonstrates environment variable parsing, validation, and startup wiring.
// No external library needed — uses only os and strconv from stdlib.
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration loaded from the environment.
type Config struct {
	Port        int
	DatabaseURL string
	LogLevel    string
	AllowedOrigins []string
}

// Load reads all required env vars, validates them, and returns a Config.
// Returns an error listing every missing or invalid variable at once (not just the first).
func Load() (*Config, error) {
	var errs []string

	port, err := requireInt("PORT", 8080)
	if err != nil {
		errs = append(errs, err.Error())
	}

	dbURL, err := requireString("DATABASE_URL")
	if err != nil {
		errs = append(errs, err.Error())
	}

	logLevel := optionalString("LOG_LEVEL", "info")
	if logLevel != "debug" && logLevel != "info" && logLevel != "warn" && logLevel != "error" {
		errs = append(errs, fmt.Sprintf("LOG_LEVEL must be one of debug|info|warn|error, got %q", logLevel))
	}

	origins := optionalStringSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"})

	if len(errs) > 0 {
		return nil, errors.New("config errors:\n  " + strings.Join(errs, "\n  "))
	}

	return &Config{
		Port:           port,
		DatabaseURL:    dbURL,
		LogLevel:       logLevel,
		AllowedOrigins: origins,
	}, nil
}

// requireString returns the value or an error if the var is unset/empty.
func requireString(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(v) == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return v, nil
}

// requireInt parses an integer env var; falls back to defaultVal if unset.
func requireInt(key string, defaultVal int) (int, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return defaultVal, nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer, got %q", key, v)
	}
	return n, nil
}

// optionalString returns the env value or the default.
func optionalString(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return defaultVal
}

// optionalStringSlice parses a comma-separated env var into a slice.
func optionalStringSlice(key string, defaultVal []string) []string {
	v, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(v) == "" {
		return defaultVal
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func main() {
	// Simulate a valid environment
	os.Setenv("PORT", "9090")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/mydb")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000, https://myapp.com")

	cfg, err := Load()
	if err != nil {
		fmt.Println("startup failed:", err)
		os.Exit(1)
	}

	fmt.Printf("config loaded: port=%d logLevel=%s origins=%v\n",
		cfg.Port, cfg.LogLevel, cfg.AllowedOrigins)

	// Simulate a missing required var
	os.Unsetenv("DATABASE_URL")
	_, err = Load()
	if err != nil {
		fmt.Println("\nmissing var error:", err)
	}
}
