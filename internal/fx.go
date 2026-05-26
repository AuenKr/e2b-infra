package internalfx

import (
	"e2b/internal/reflection"
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
		fx.Annotate(
			reflection.NewReflectionRouteV1,
			fx.ResultTags(`group:"grpc-routes"`),
		),
		fx.Annotate(
			reflection.NewReflectionRouteV1Alpha,
			fx.ResultTags(`group:"grpc-routes"`),
		),
	),
	fx.Provide(server.NewServerMux),
)
