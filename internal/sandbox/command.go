package sandbox

import (
	"bytes"
	"context"
	"strings"

	sandboxv1 "e2b/gen/sandbox/v1"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	utilexec "k8s.io/client-go/util/exec"
)

// SendCommand implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) SendCommand(ctx context.Context, req *sandboxv1.SendCommandRequest) (*sandboxv1.SendCommandResponse, error) {
	// Send Exec cmd to pod
	k8core := s.K8sClient.CoreV1()

	cmd := strings.Fields(req.GetCommand())

	execReq := k8core.RESTClient().
		Post().
		Namespace(s.Config.K8sNamespace).
		Resource("pods").
		Name(req.GetId()).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
			Command: cmd,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(s.K8sConfig, "POST", execReq.URL())
	if err != nil {
		return nil, err
	}

	var stdout, stderr bytes.Buffer
	if err := executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	}); err != nil {
		if _, ok := err.(utilexec.CodeExitError); !ok {
			s.Logger.Error("failed to exec command", zap.Error(err))
			return nil, err
		}
	}

	return &sandboxv1.SendCommandResponse{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, nil
}
