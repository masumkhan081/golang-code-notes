package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port        int
	DatabaseURL string
	LogLevel    string
}

func Load() (*Config, error) {
	var errs []string

	port, err := parseInt("PORT", 8080)
	if err != nil {
		errs = append(errs, err.Error())
	}

	dbURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok || strings.TrimSpace(dbURL) == "" {
		errs = append(errs, "DATABASE_URL is required")
	}

	logLevel := optionalString("LOG_LEVEL", "info")

	if len(errs) > 0 {
		return nil, errors.New("config errors: " + strings.Join(errs, "; "))
	}

	return &Config{
		Port:        port,
		DatabaseURL: strings.TrimSpace(dbURL),
		LogLevel:    logLevel,
	}, nil
}

func parseInt(key string, def int) (int, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer, got %q", key, v)
	}
	return n, nil
}

func optionalString(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return def
}
