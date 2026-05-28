package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Mode          string
	Port          int
	K8sConfigPath string
	K8sNamespace  string
}

const (
	CPU_MIN_DEFAULT    = "50m"
	CPU_MAX_DEFAULT    = "4000m"
	MEMORY_MIN_DEFAULT = "1Mi"
	MEMORY_MAX_DEFAULT = "8Gi"
)

func NewConfig() Config {
	_ = godotenv.Load()
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(fmt.Sprintln("Invlaid Port", err))
	}

	Mode := getEnv("MODE", "dev")
	K8sConfigPath := getEnv("K8S_CONFIG_PATH", "")
	K8sNamespace := getEnv("K8S_NAMESPACE", "default")

	return Config{
		Port:          port,
		Mode:          Mode,
		K8sConfigPath: K8sConfigPath,
		K8sNamespace:  K8sNamespace,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
