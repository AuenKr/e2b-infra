package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port int
}

func NewConfig() Config {
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(fmt.Sprintln("Invlaid Port", err))
	}

	return Config{
		Port: port,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
