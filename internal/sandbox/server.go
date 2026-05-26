// Package sandbox: Implements the sandbox.v1.SandboxService service.
package sandbox

import (
	"context"

	sandboxv1 "e2b/gen/sandbox/v1"
	"e2b/gen/sandbox/v1/sandboxv1connect"
	"e2b/pkg/config"
	"e2b/pkg/server"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"go.uber.org/fx"
	"k8s.io/client-go/kubernetes"
)

type SandboxServerParams struct {
	fx.In
	Config config.Config
}

func NewSandboxServer(in SandboxServerParams) *SandboxServer {
	return &SandboxServer{
		config: in.Config,
	}
}

type SandboxServerRouteParams struct {
	fx.In
	Server    *SandboxServer
	K8sClient *kubernetes.Clientset
}

func NewSandboxServerRoute(in SandboxServerRouteParams) server.Route {
	p, h := sandboxv1connect.NewSandboxServiceHandler(
		in.Server,
		connect.WithInterceptors(validate.NewInterceptor()),
	)
	return server.Route{Path: p, Handler: h}
}

// For compile time error
var _ sandboxv1connect.SandboxServiceHandler = (*SandboxServer)(nil)

type SandboxServer struct {
	config config.Config
}

// StartSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) StartSandbox(context.Context, *sandboxv1.StartSandboxRequest) (*sandboxv1.StartSandboxResponse, error) {
	// Start the pod with given image
	panic("unimplemented")
}

// StopSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) StopSandbox(context.Context, *sandboxv1.StopSandboxRequest) (*sandboxv1.StopSandboxResponse, error) {
	// Stop the pod
	panic("unimplemented")
}

// GetSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) GetSandbox(context.Context, *sandboxv1.GetSandboxRequest) (*sandboxv1.GetSandboxResponse, error) {
	// Get pod details
	panic("unimplemented")
}

// ListSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) ListSandbox(context.Context, *sandboxv1.ListSandboxRequest) (*sandboxv1.ListSandboxResponse, error) {
	// Get user pods
	panic("unimplemented")
}

// SendCommand implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) SendCommand(context.Context, *sandboxv1.SendCommandRequest) (*sandboxv1.SendCommandResponse, error) {
	// Send Exec cmd to pod
	panic("unimplemented")
}

// OpenPort implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) OpenPort(context.Context, *sandboxv1.OpenPortRequest) (*sandboxv1.OpenPortResponse, error) {
	// Create a service to pod with that open port
	panic("unimplemented")
}

// ClosePort implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) ClosePort(context.Context, *sandboxv1.ClosePortRequest) (*sandboxv1.ClosePortResponse, error) {
	// Close the open port in the service for that pod
	panic("unimplemented")
}

// ListOpenPort implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) ListOpenPort(context.Context, *sandboxv1.ListOpenPortRequest) (*sandboxv1.ListOpenPortResponse, error) {
	// Read the service for that pod and get its open ports
	panic("unimplemented")
}
