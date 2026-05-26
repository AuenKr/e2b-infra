// Package internalfx: Internal module for the server.
package internalfx

import (
	"e2b/internal/sandbox"
	"e2b/internal/server"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"internal",
	fx.Provide(sandbox.NewSandboxServer),
	fx.Provide(
		fx.Annotate(
			sandbox.NewSandboxServerRoute,
			fx.ResultTags(`group:"grpc-routes"`),
		),
	),
	fx.Provide(server.NewHTTPMux),
	fx.Decorate(server.RegisterRoutes),
)
