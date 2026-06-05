package internalfx

import (
	"e2b/internal/clients"
	"e2b/internal/port_manger"
	"e2b/internal/reflection"
	"e2b/internal/sandbox"
	"e2b/internal/server"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"internal",
	clients.Module,
	server.Module,
	reflection.Module,

	// Service
	sandbox.Module,
	port_manger.Module,
)
