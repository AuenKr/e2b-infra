package port_manger

import (
	"context"
	"errors"

	commonv1 "e2b/gen/common/v1"
	portmangerv1 "e2b/gen/port_manger/v1"
	"e2b/internal/sandbox"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// OpenPort implements [port_mangerv1connect.PortMangerServiceHandler].
func (s *PortMangerServer) OpenPort(ctx context.Context, req *portmangerv1.OpenPortRequest) (*portmangerv1.OpenPortResponse, error) {
	// Create a service to pod with that open port
	k8core := s.K8sClient.CoreV1()

	service, err := k8core.Services(s.Config.K8sNamespace).Get(ctx, req.GetId(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	for _, port := range service.Spec.Ports {
		if port.Port == int32(req.Port.PortNumber) && port.Protocol == sandbox.K8sProtocolAdapter(req.Port.Protocol) {
			return nil, errors.New("already opened")
		}
	}

	newPort := corev1.ServicePort{
		Name:     sandbox.GetPortName(req.Port),
		Protocol: sandbox.K8sProtocolAdapter(req.Port.Protocol),
		Port:     int32(req.Port.PortNumber),
	}

	service.Spec.Ports = append(service.Spec.Ports, newPort)
	_, err = k8core.Services(s.Config.K8sNamespace).Update(ctx, service, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	httpRoutes, err := s.GatewayClient.GatewayV1beta1().HTTPRoutes(s.Config.K8sNamespace).Get(ctx, req.GetId(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	portHostname := sandbox.GetHTTPRoute(req.Port, s.Config.Domain)
	for _, hostname := range httpRoutes.Spec.Hostnames {
		if string(hostname) == portHostname {
			return nil, errors.New("already registered http route")
		}
	}
	httpRoutes.Spec.Hostnames = append(httpRoutes.Spec.Hostnames, gatewayapiv1.Hostname(portHostname))

	ruleName := gatewayapiv1.SectionName(sandbox.GetPortName(req.Port))
	for _, rule := range httpRoutes.Spec.Rules {
		if rule.Name != nil && *rule.Name == ruleName {
			return nil, errors.New("already defined http route rule")
		}
	}

	portNumber := gatewayapiv1.PortNumber(int32(req.Port.PortNumber))
	rules := gatewayapiv1.HTTPRouteRule{
		Name: &ruleName,
		Matches: []gatewayapiv1.HTTPRouteMatch{
			{
				Headers: []gatewayapiv1.HTTPHeaderMatch{
					{
						Name:  "Host",
						Value: portHostname,
					},
				},
			},
		},
		BackendRefs: []gatewayapiv1.HTTPBackendRef{
			{
				BackendRef: gatewayapiv1.BackendRef{
					BackendObjectReference: gatewayapiv1.BackendObjectReference{
						Name: gatewayapiv1.ObjectName(req.Id),
						Port: &portNumber,
					},
				},
			},
		},
	}

	httpRoutes.Spec.Rules = append(httpRoutes.Spec.Rules, rules)

	_, err = s.GatewayClient.GatewayV1beta1().HTTPRoutes(s.Config.K8sNamespace).Update(ctx, httpRoutes, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return &portmangerv1.OpenPortResponse{
		Port: &commonv1.PortInfo{
			PortNumber: req.Port.PortNumber,
			Protocol:   req.Port.Protocol,
			Hostname:   sandbox.GetHTTPRoute(req.Port, s.Config.Domain),
		},
	}, nil
}

// ClosePort implements [port_mangerv1connect.PortMangerServiceHandler].
func (s *PortMangerServer) ClosePort(ctx context.Context, req *portmangerv1.ClosePortRequest) (*portmangerv1.ClosePortResponse, error) {
	// Close the open port in the service for that pod
	// Create a service to pod with that open port
	k8core := s.K8sClient.CoreV1()

	service, err := k8core.Services(s.Config.K8sNamespace).Get(ctx, req.GetId(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	newPorts := make([]corev1.ServicePort, 0, len(service.Spec.Ports))
	for _, port := range service.Spec.Ports {
		if port.Port == int32(req.Port.PortNumber) && port.Protocol == sandbox.K8sProtocolAdapter(req.Port.Protocol) {
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

	httpRoutes, err := s.GatewayClient.GatewayV1beta1().HTTPRoutes(s.Config.K8sNamespace).Get(ctx, req.GetId(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	hostnames := make([]gatewayapiv1.Hostname, 0, len(httpRoutes.Spec.Hostnames))
	route := sandbox.GetHTTPRoute(req.Port, s.Config.Domain)
	for _, host := range httpRoutes.Spec.Hostnames {
		if string(host) == route {
			continue
		}
		hostnames = append(hostnames, host)
	}
	if len(hostnames) == len(httpRoutes.Spec.Hostnames) {
		return nil, errors.New("HTTP hostnames does not exist")
	}

	rules := make([]gatewayapiv1.HTTPRouteRule, 0, len(httpRoutes.Spec.Rules))
	ruleName := gatewayapiv1.SectionName(sandbox.GetPortName(req.Port))
	for _, rule := range httpRoutes.Spec.Rules {
		if rule.Name != nil && *rule.Name == ruleName {
			continue
		}
		rules = append(rules, rule)
	}
	if len(rules) == len(httpRoutes.Spec.Rules) {
		return nil, errors.New("HTTP route does not exist")
	}

	httpRoutes.Spec.Hostnames = hostnames
	httpRoutes.Spec.Rules = rules

	_, err = s.GatewayClient.GatewayV1beta1().HTTPRoutes(s.Config.K8sNamespace).Update(ctx, httpRoutes, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return &portmangerv1.ClosePortResponse{
		Port: &commonv1.PortInfo{
			PortNumber: req.Port.PortNumber,
			Protocol:   req.Port.Protocol,
			Hostname:   sandbox.GetHTTPRoute(req.Port, s.Config.Domain),
		},
	}, nil
}

// ListOpenPort implements [port_mangerv1connect.PortMangerServiceHandler].
func (s *PortMangerServer) ListOpenPort(ctx context.Context, req *portmangerv1.ListOpenPortRequest) (*portmangerv1.ListOpenPortResponse, error) {
	// Read the service for that pod and get its open ports
	k8core := s.K8sClient.CoreV1()

	k8sPod, err := k8core.Services(s.Config.K8sNamespace).Get(ctx, req.GetId(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	ports := make([]*commonv1.PortInfo, len(k8sPod.Spec.Ports))
	for i, port := range k8sPod.Spec.Ports {
		port := commonv1.PortInfo{
			PortNumber: uint32(port.Port),
			Protocol:   sandbox.SandboxProtocolAdapter(port.Protocol),
		}
		hostname := sandbox.GetHTTPRoute(&port, s.Config.Domain)
		port.Hostname = hostname
		ports[i] = &port
	}
	return &portmangerv1.ListOpenPortResponse{
		Ports: ports,
	}, nil
}
