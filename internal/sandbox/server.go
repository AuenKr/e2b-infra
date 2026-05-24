// Package sandbox: Implements the sandbox.v1.SandboxService service.
package sandbox

import (
	"context"

	sandboxv1 "e2b/gen/sandbox/v1"
	"e2b/gen/sandbox/v1/sandboxv1connect"
	"e2b/pkg/config"

	"go.uber.org/fx"
)

type SandboxServerParams struct {
	fx.In
	config config.Config
}

func NewSandboxServer(in SandboxServerParams) *SandboxServer {
	return &SandboxServer{
		config: in.config,
	}
}

// For compile time error
var _ sandboxv1connect.SandboxServiceHandler = (*SandboxServer)(nil)

type SandboxServer struct {
	config config.Config
}

// StartSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) StartSandbox(context.Context, *sandboxv1.StartSandboxRequest) (*sandboxv1.StartSandboxResponse, error) {
	panic("unimplemented")
}

// StopSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) StopSandbox(context.Context, *sandboxv1.StopSandboxRequest) (*sandboxv1.StopSandboxResponse, error) {
	panic("unimplemented")
}

// GetSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) GetSandbox(context.Context, *sandboxv1.GetSandboxRequest) (*sandboxv1.GetSandboxResponse, error) {
	panic("unimplemented")
}

// ListSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) ListSandbox(context.Context, *sandboxv1.ListSandboxRequest) (*sandboxv1.ListSandboxResponse, error) {
	panic("unimplemented")
}

// SendCommand implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) SendCommand(context.Context, *sandboxv1.SendCommandRequest) (*sandboxv1.SendCommandResponse, error) {
	panic("unimplemented")
}

// OpenPort implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) OpenPort(context.Context, *sandboxv1.OpenPortRequest) (*sandboxv1.OpenPortResponse, error) {
	panic("unimplemented")
}

// ClosePort implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) ClosePort(context.Context, *sandboxv1.ClosePortRequest) (*sandboxv1.ClosePortResponse, error) {
	panic("unimplemented")
}

// ListOpenPort implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) ListOpenPort(context.Context, *sandboxv1.ListOpenPortRequest) (*sandboxv1.ListOpenPortResponse, error) {
	panic("unimplemented")
}
