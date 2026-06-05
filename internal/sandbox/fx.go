package sandbox

import "go.uber.org/fx"

var Module = fx.Module(
	"sandbox",
	fx.Provide(NewSandboxServer),
	fx.Provide(
		fx.Annotate(
			NewSandboxServerRoute,
			fx.ResultTags(`group:"grpc-routes"`),
		),
	),
)
