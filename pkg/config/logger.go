package config

import "go.uber.org/zap"

func NewLogger(config Config) *zap.Logger {
	var logger *zap.Logger
	if config.Mode == "dev" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	return logger
}
