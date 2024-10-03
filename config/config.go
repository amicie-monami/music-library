package config

import (
	"flag"
	"log"
	"log/slog"
	"strconv"

	"github.com/joho/godotenv"
)

var envFilePath *string
var appEnv *string

func init() {
	appEnv = flag.String("env", "stage", "application environment, can take one value from [stage, test] ")
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
	LogLevel string
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

	dbSourceEnvVar := "DATABASE_SOURCE"
	if *appEnv == "test" {
		dbSourceEnvVar = "TEST_DATABASE_SOURCE"
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
			Source: env[dbSourceEnvVar],
		},
		LogLevel: env["LOG_LEVEL"],
	}
}

func mustParseDigit(raw string) int {
	num, err := strconv.Atoi(raw)
	if err != nil {
		log.Fatalf("failed to parse str=%s", raw)
	}
	return num
}

func ConfigureSlogLogger(logLevel string) {
	if logLevel == "debug" {
		slog.SetLogLoggerLevel(slog.LevelDebug)

	} else if logLevel == "info" {
		slog.SetLogLoggerLevel(slog.LevelInfo)

	} else if logLevel == "error" {
		slog.SetLogLoggerLevel(slog.LevelError)

	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
}
