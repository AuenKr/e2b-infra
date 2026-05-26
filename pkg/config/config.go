package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Mode string
	Port int
}

func NewConfig() Config {
	_ = godotenv.Load()
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(fmt.Sprintln("Invlaid Port", err))
	}

	Mode := getEnv("MODE", "dev")

	return Config{
		Port: port,
		Mode: Mode,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
