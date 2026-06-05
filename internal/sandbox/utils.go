package sandbox

import (
	"fmt"
	"strings"

	commonv1 "e2b/gen/common/v1"

	corev1 "k8s.io/api/core/v1"
)

func PodPhaseToSandboxStateAdapter(status corev1.PodPhase) commonv1.SandboxStatus {
	var sandboxStatus commonv1.SandboxStatus
	switch status {
	case corev1.PodPending:
		sandboxStatus = commonv1.SandboxStatus_SANDBOX_STATUS_STARTING
	case corev1.PodRunning:
		sandboxStatus = commonv1.SandboxStatus_SANDBOX_STATUS_RUNNING
	case corev1.PodSucceeded:
		sandboxStatus = commonv1.SandboxStatus_SANDBOX_STATUS_STOPPED
	default:
		// corev1.PodFailed, corev1.PodUnknown
		sandboxStatus = commonv1.SandboxStatus_SANDBOX_STATUS_ERROR
	}
	return sandboxStatus
}

func K8sProtocolAdapter(protocol commonv1.Protocol) corev1.Protocol {
	var coreProtocol corev1.Protocol
	switch protocol {
	case commonv1.Protocol_PROTOCOL_TCP:
		coreProtocol = corev1.ProtocolTCP
	case commonv1.Protocol_PROTOCOL_UDP:
		coreProtocol = corev1.ProtocolUDP
	case commonv1.Protocol_PROTOCOL_SCTP:
		coreProtocol = corev1.ProtocolSCTP
	}
	return coreProtocol
}

func SandboxProtocolAdapter(protocol corev1.Protocol) commonv1.Protocol {
	var sandboxProtocol commonv1.Protocol
	switch protocol {
	case corev1.ProtocolTCP:
		sandboxProtocol = commonv1.Protocol_PROTOCOL_TCP
	case corev1.ProtocolUDP:
		sandboxProtocol = commonv1.Protocol_PROTOCOL_UDP
	case corev1.ProtocolSCTP:
		sandboxProtocol = commonv1.Protocol_PROTOCOL_SCTP
	}
	return sandboxProtocol
}

func GetPortName(port *commonv1.PortInfo) string {
	name := strings.ToLower(port.Protocol.String())
	formattedName := strings.ReplaceAll(name, "_", "-")
	return strings.ToLower(fmt.Sprintf("%s-%d", formattedName, port.PortNumber))
}

func GetHTTPRoute(port *commonv1.PortInfo, domain string) string {
	return fmt.Sprintf("%s.%s", GetPortName(port), domain)
}
