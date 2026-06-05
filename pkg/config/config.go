package config

import (
	"fmt"
	"os"
	"strconv"

	commonv1 "e2b/gen/common/v1"

	"github.com/joho/godotenv"
)

type Config struct {
	Mode             string
	ReflectionEnable bool

	Port          int
	K8sConfigPath string
	K8sNamespace  string

	K8sGateway          string
	K8sGatewayNamespace string

	Domain string
}

const (
	CPU_MIN_DEFAULT    = "50m"
	CPU_MAX_DEFAULT    = "4000m"
	MEMORY_MIN_DEFAULT = "1Mi"
	MEMORY_MAX_DEFAULT = "8Gi"
)

const (
	INITIAL_PORT          = 69
	INITIAL_PORT_NAME     = "intial-port"
	INITIAL_PORT_PROTOCOL = commonv1.Protocol_PROTOCOL_TCP
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

	K8sGateway := getEnv("K8S_GATEWAY", "")
	if K8sGateway == "" {
		panic("K8sGateway value required")
	}
	K8sGatewayNamespace := getEnv("K8S_GATEWAY_NAMESPACE", "default")
	Domain := getEnv("DOMAIN_NAME", "")
	if Domain == "" {
		panic("Domain Value Required")
	}

	reflectionEnable := true
	reflectionEnvEnable := getEnv("REFLECTION_ENABLE", "true")
	if reflectionEnvEnable == "false" {
		reflectionEnable = false
	}

	return Config{
		Port:                port,
		ReflectionEnable:    reflectionEnable,
		Mode:                Mode,
		K8sConfigPath:       K8sConfigPath,
		K8sNamespace:        K8sNamespace,
		K8sGateway:          K8sGateway,
		K8sGatewayNamespace: K8sGatewayNamespace,
		Domain:              Domain,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
