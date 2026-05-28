package sandbox

import (
	"context"
	"errors"

	sandboxv1 "e2b/gen/sandbox/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenPort implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) OpenPort(ctx context.Context, req *sandboxv1.OpenPortRequest) (*sandboxv1.OpenPortResponse, error) {
	// Create a service to pod with that open port
	k8core := s.K8sClient.CoreV1()

	service, err := k8core.Services(s.Config.K8sNamespace).Get(ctx, req.GetId(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	for _, port := range service.Spec.Ports {
		if port.Port == int32(req.Port.PortNumber) && port.Protocol == K8sProtocolAdapter(req.Port.Protocol) {
			return nil, errors.New("already opened")
		}
	}

	newPort := corev1.ServicePort{
		Name:     GetPortName(req.Port),
		Protocol: K8sProtocolAdapter(req.Port.Protocol),
		Port:     int32(req.Port.PortNumber),
	}

	service.Spec.Ports = append(service.Spec.Ports, newPort)
	_, err = k8core.Services(s.Config.K8sNamespace).Update(ctx, service, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return &sandboxv1.OpenPortResponse{
		Port: &sandboxv1.PortInfo{
			PortNumber: req.Port.PortNumber,
			Protocol:   req.Port.Protocol,
			// TODO: Attach some url
			// Hostname:   "",
		},
	}, nil
}

// ClosePort implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) ClosePort(ctx context.Context, req *sandboxv1.ClosePortRequest) (*sandboxv1.ClosePortResponse, error) {
	// Close the open port in the service for that pod
	// Create a service to pod with that open port
	k8core := s.K8sClient.CoreV1()

	service, err := k8core.Services(s.Config.K8sNamespace).Get(ctx, req.GetId(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	newPorts := make([]corev1.ServicePort, 0, len(service.Spec.Ports))
	for _, port := range service.Spec.Ports {
		if port.Port == int32(req.Port.PortNumber) && port.Protocol == K8sProtocolAdapter(req.Port.Protocol) {
			continue
		}
		newPorts = append(newPorts, port)
	}
	if len(newPorts) == len(service.Spec.Ports) {
		return nil, errors.New("port is not opened")
	}

	service.Spec.Ports = newPorts
	_, err = k8core.Services(s.Config.K8sNamespace).Update(ctx, service, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return &sandboxv1.ClosePortResponse{
		Port: &sandboxv1.PortInfo{
			PortNumber: req.Port.PortNumber,
			Protocol:   req.Port.Protocol,
			// TODO: Attach some url
			// Hostname:   "",
		},
	}, nil
}

// ListOpenPort implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) ListOpenPort(ctx context.Context, req *sandboxv1.ListOpenPortRequest) (*sandboxv1.ListOpenPortResponse, error) {
	// Read the service for that pod and get its open ports
	k8core := s.K8sClient.CoreV1()

	k8sPod, err := k8core.Services(s.Config.K8sNamespace).Get(ctx, req.GetId(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	ports := make([]*sandboxv1.PortInfo, len(k8sPod.Spec.Ports))
	for i, port := range k8sPod.Spec.Ports {
		ports[i] = &sandboxv1.PortInfo{
			PortNumber: uint32(port.Port),
			Protocol:   SandboxProtocolAdapter(port.Protocol),
			// Hostname: port.,
		}
	}
	return &sandboxv1.ListOpenPortResponse{
		Ports: ports,
	}, nil
}
