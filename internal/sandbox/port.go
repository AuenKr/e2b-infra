package sandbox

import (
	"context"

	sandboxv1 "e2b/gen/sandbox/v1"
)

// OpenPort implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) OpenPort(ctx context.Context, req *sandboxv1.OpenPortRequest) (*sandboxv1.OpenPortResponse, error) {
	// Create a service to pod with that open port
	//
	panic("unimplemented")
}

// ClosePort implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) ClosePort(ctx context.Context, req *sandboxv1.ClosePortRequest) (*sandboxv1.ClosePortResponse, error) {
	// Close the open port in the service for that pod
	panic("unimplemented")
}

// ListOpenPort implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) ListOpenPort(ctx context.Context, req *sandboxv1.ListOpenPortRequest) (*sandboxv1.ListOpenPortResponse, error) {
	// Read the service for that pod and get its open ports
	panic("unimplemented")
}
