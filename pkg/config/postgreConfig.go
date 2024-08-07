package config

import (
	"fmt"
	"os"
)

type PostgreConfig struct {
	DBName   string
	User     string
	Password string
	Port     string
	Host     string
	Schema   string
}

func InitPostgreConfig() *PostgreConfig {
	config := &PostgreConfig{
		Host:     getEnvOrPanic("DB_HOST"),
		Port:     getEnvOrPanic("DB_PORT"),
		User:     getEnvOrPanic("DB_USER"),
		Password: getEnvOrPanic("DB_PASS"),
		DBName:   getEnvOrPanic("DB_NAME"),
		Schema:   getEnvOrPanic("DB_SCHEMA"),
	}
	return config
}

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s is not set", key))
	}
	return value
}
