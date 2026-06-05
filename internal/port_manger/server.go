package port_manger

import (
	"e2b/gen/port_manger/v1/port_mangerv1connect"
	"e2b/internal/middleware"
	"e2b/internal/sandbox"
	"e2b/pkg/server"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PortMangerServer struct {
	*sandbox.SandboxServer
}

func NewPortMangerServer(server *sandbox.SandboxServer) *PortMangerServer {
	return &PortMangerServer{SandboxServer: server}
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
