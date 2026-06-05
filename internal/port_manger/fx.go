package port_manger

import "go.uber.org/fx"

var Module = fx.Module(
	"port-manager",
	fx.Provide(NewPortMangerServer),
	fx.Provide(
		fx.Annotate(
			NewPortMangerServerRoute,
			fx.ResultTags(`group:"grpc-routes"`),
		),
	),
)
