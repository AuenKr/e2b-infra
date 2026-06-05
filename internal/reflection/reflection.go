package reflection

import (
	"e2b/pkg/server"

	"connectrpc.com/grpcreflect"
	"go.uber.org/fx"
)

type ReflectionRouteV1Params struct {
	fx.In
	Routes []server.Route `group:"grpc-routes"`
}

func NewReflectionRouteV1(in ReflectionRouteV1Params) server.Route {
	serviceNames := make([]string, len(in.Routes))
	for i, r := range in.Routes {
		serviceNames[i] = r.ServiceName
	}

	p, h := grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(
			serviceNames...,
		),
	)
	return server.Route{Path: p, Handler: h}
}

type ReflectionRouteV1AlphaParams struct {
	fx.In
	Routes []server.Route `group:"grpc-routes"`
}

func NewReflectionRouteV1Alpha(in ReflectionRouteV1AlphaParams) server.Route {
	serviceNames := make([]string, len(in.Routes))
	for i, r := range in.Routes {
		serviceNames[i] = r.ServiceName
	}

	p, h := grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(
			serviceNames...,
		),
	)
	return server.Route{Path: p, Handler: h}
}
