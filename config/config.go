package config

import (
	"flag"
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

var envFilePath *string

func init() {
	envFilePath = flag.String("config", ".env", "path to the .env configuration file")
}

// ServerConfig stores the configuration of the http server
type ServerConfig struct {
	Addr           string
	ReadTimeout    int
	WriteTimeout   int
	IdleTimeout    int
	MaxHeaderBytes int
}

type DatabaseConfig struct {
	Source string
}

// Config stores the configuration of the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// MustLoadFromEnv loads config from .env file
// To change standart path use flag --config
func MustLoadFromEnv() *Config {
	if err := godotenv.Load(*envFilePath); err != nil {
		log.Fatal("failed to load .env file", err)
	}

	env, err := godotenv.Read(*envFilePath)
	if err != nil {
		log.Fatal("failed to parse .env file", err)
	}

	return &Config{
		Server: ServerConfig{
			Addr:           env["SERVER_ADDR"],
			ReadTimeout:    mustParseDigit(env["SERVER_READ_TIMEOUT"]),
			WriteTimeout:   mustParseDigit(env["SERVER_WRITE_TIMEOUT"]),
			IdleTimeout:    mustParseDigit(env["SERVER_IDLE_TIMEOUT"]),
			MaxHeaderBytes: mustParseDigit(env["SERVER_MAX_HEADER_BYTES"]),
		},
		Database: DatabaseConfig{
			Source: env["DATABASE_SOURCE"],
		},
	}
}

func mustParseDigit(raw string) int {
	num, err := strconv.Atoi(raw)
	if err != nil {
		log.Fatalf("failed to parse str=%s", raw)
	}
	return num
}
