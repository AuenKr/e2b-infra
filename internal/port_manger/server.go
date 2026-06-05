package port_manger

import (
	"e2b/gen/port_manger/v1/port_mangerv1connect"
	"e2b/internal/middleware"
	"e2b/pkg/config"
	"e2b/pkg/server"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	gatewayclientset "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

type PortMangerServer struct {
	Config        config.Config
	K8sClient     *kubernetes.Clientset
	GatewayClient *gatewayclientset.Clientset
	Logger        *zap.Logger
}

type PortMangerServerParams struct {
	fx.In
	Config        config.Config
	K8sClient     *kubernetes.Clientset
	GatewayClient *gatewayclientset.Clientset
	Logger        *zap.Logger
}

func NewPortMangerServer(in PortMangerServerParams) *PortMangerServer {
	return &PortMangerServer{
		Config:        in.Config,
		K8sClient:     in.K8sClient,
		GatewayClient: in.GatewayClient,
		Logger:        in.Logger,
	}
}

var _ port_mangerv1connect.PortMangerServiceHandler = (*PortMangerServer)(nil)

type PortMangerServerRouteParams struct {
	fx.In
	Server *PortMangerServer
	Logger *zap.Logger
}

func NewPortMangerServerRoute(in PortMangerServerRouteParams) server.Route {
	p, h := port_mangerv1connect.NewPortMangerServiceHandler(
		in.Server,
		connect.WithInterceptors(
			middleware.NewLoggerInterceptor(in.Logger),
			validate.NewInterceptor(),
		),
	)
	return server.Route{
		ServiceName: port_mangerv1connect.PortMangerServiceName,
		Path:        p,
		Handler:     h,
	}
}
