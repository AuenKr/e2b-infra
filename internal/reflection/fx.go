package reflection

import "go.uber.org/fx"

var Module = fx.Module(
	"reflection",
	fx.Provide(
		fx.Annotate(
			NewReflectionRouteV1,
			fx.ResultTags(`group:"reflection"`),
		),
		fx.Annotate(
			NewReflectionRouteV1Alpha,
			fx.ResultTags(`group:"reflection"`),
		),
	),
)
