package reflection

import (
	"e2b/gen/sandbox/v1/sandboxv1connect"
	"e2b/pkg/server"

	"connectrpc.com/grpcreflect"
)

func NewReflectionRouteV1() server.Route {
	p, h := grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(sandboxv1connect.SandboxServiceName),
	)
	return server.Route{Path: p, Handler: h}
}

func NewReflectionRouteV1Alpha() server.Route {
	p, h := grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(sandboxv1connect.SandboxServiceName),
	)
	return server.Route{Path: p, Handler: h}
}
