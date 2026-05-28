// Package sandbox: Implements the sandbox.v1.SandboxService service.
package sandbox

import (
	"e2b/gen/sandbox/v1/sandboxv1connect"
	"e2b/internal/middleware"
	"e2b/pkg/config"
	"e2b/pkg/server"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"k8s.io/client-go/kubernetes"
)

type SandboxServerParams struct {
	fx.In
	Config    config.Config
	K8sClient *kubernetes.Clientset
	Logger    *zap.Logger
}

func NewSandboxServer(in SandboxServerParams) *SandboxServer {
	return &SandboxServer{
		Config:    in.Config,
		K8sClient: in.K8sClient,
		Logger:    in.Logger,
	}
}

type SandboxServerRouteParams struct {
	fx.In
	Server *SandboxServer
	Logger *zap.Logger
}

func NewSandboxServerRoute(in SandboxServerRouteParams) server.Route {
	p, h := sandboxv1connect.NewSandboxServiceHandler(
		in.Server,
		connect.WithInterceptors(
			middleware.NewLoggerInterceptor(in.Logger),
			validate.NewInterceptor(),
		),
	)
	return server.Route{Path: p, Handler: h}
}

// For compile time error
var _ sandboxv1connect.SandboxServiceHandler = (*SandboxServer)(nil)

type SandboxServer struct {
	Config    config.Config
	K8sClient *kubernetes.Clientset
	Logger    *zap.Logger
}
