package sandbox

import (
	"errors"

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

func protocolAdapter(protocol sandboxv1.Protocol) (corev1.Protocol, error) {
	var coreProtocol corev1.Protocol
	switch protocol {
	case sandboxv1.Protocol_PROTOCOL_TCP:
		coreProtocol = corev1.ProtocolTCP
	case sandboxv1.Protocol_PROTOCOL_UDP:
		coreProtocol = corev1.ProtocolUDP
	case sandboxv1.Protocol_PROTOCOL_SCTP:
		coreProtocol = corev1.ProtocolSCTP
	default:
		return coreProtocol, errors.New("invalid protocol")
	}
	return coreProtocol, nil
}
