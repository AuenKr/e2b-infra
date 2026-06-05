package clients

import "go.uber.org/fx"

var Module = fx.Module(
	"clients",
	fx.Provide(NewK8sRESTConfig),
	fx.Provide(NewK8sClusterClient),
	fx.Provide(NewK8sMetricClient),
	fx.Provide(NewK8sGatewayClient),
)
