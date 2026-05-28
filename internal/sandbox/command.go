package sandbox

import (
	"context"

	sandboxv1 "e2b/gen/sandbox/v1"
)

// SendCommand implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) SendCommand(ctx context.Context, req *sandboxv1.SendCommandRequest) (*sandboxv1.SendCommandResponse, error) {
	// Send Exec cmd to pod
	// k8core := s.K8sClient.CoreV1()

	panic("unimplemented")
}
