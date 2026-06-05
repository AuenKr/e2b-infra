package internalfx

import (
	"e2b/internal/port_manger"
	"e2b/internal/reflection"
	"e2b/internal/sandbox"
	"e2b/internal/server"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"internal",
	fx.Provide(sandbox.NewSandboxServer),
	fx.Provide(port_manger.NewPortMangerServer),
	fx.Provide(
		fx.Annotate(
			sandbox.NewSandboxServerRoute,
			fx.ResultTags(`group:"grpc-routes"`),
		),
		fx.Annotate(
			port_manger.NewPortMangerServerRoute,
			fx.ResultTags(`group:"grpc-routes"`),
		),
		fx.Annotate(
			reflection.NewReflectionRouteV1,
			fx.ResultTags(`group:"reflection"`),
		),
		fx.Annotate(
			reflection.NewReflectionRouteV1Alpha,
			fx.ResultTags(`group:"reflection"`),
		),
	),
	fx.Provide(server.NewServerMux),
)
