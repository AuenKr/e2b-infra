package sandbox

import (
	"fmt"
	"strings"

	sandboxv1 "e2b/gen/sandbox/v1"

	corev1 "k8s.io/api/core/v1"
)

func PodPhaseToSandboxStateAdapter(status corev1.PodPhase) sandboxv1.SandboxStatus {
	var sandboxStatus sandboxv1.SandboxStatus
	switch status {
	case corev1.PodPending:
		sandboxStatus = sandboxv1.SandboxStatus_SANDBOX_STATUS_STARTING
	case corev1.PodRunning:
		sandboxStatus = sandboxv1.SandboxStatus_SANDBOX_STATUS_RUNNING
	case corev1.PodSucceeded:
		sandboxStatus = sandboxv1.SandboxStatus_SANDBOX_STATUS_STOPPED
	default:
		// corev1.PodFailed, corev1.PodUnknown
		sandboxStatus = sandboxv1.SandboxStatus_SANDBOX_STATUS_ERROR
	}
	return sandboxStatus
}

func K8sProtocolAdapter(protocol sandboxv1.Protocol) corev1.Protocol {
	var coreProtocol corev1.Protocol
	switch protocol {
	case sandboxv1.Protocol_PROTOCOL_TCP:
		coreProtocol = corev1.ProtocolTCP
	case sandboxv1.Protocol_PROTOCOL_UDP:
		coreProtocol = corev1.ProtocolUDP
	case sandboxv1.Protocol_PROTOCOL_SCTP:
		coreProtocol = corev1.ProtocolSCTP
	}
	return coreProtocol
}

func SandboxProtocolAdapter(protocol corev1.Protocol) sandboxv1.Protocol {
	var sandboxProtocol sandboxv1.Protocol
	switch protocol {
	case corev1.ProtocolTCP:
		sandboxProtocol = sandboxv1.Protocol_PROTOCOL_TCP
	case corev1.ProtocolUDP:
		sandboxProtocol = sandboxv1.Protocol_PROTOCOL_UDP
	case corev1.ProtocolSCTP:
		sandboxProtocol = sandboxv1.Protocol_PROTOCOL_SCTP
	}
	return sandboxProtocol
}

func GetPortName(port *sandboxv1.PortInfo) string {
	name := strings.ToLower(port.Protocol.String())
	formattedName := strings.ReplaceAll(name, "_", "-")
	return strings.ToLower(fmt.Sprintf("%s-%d", formattedName, port.PortNumber))
}
