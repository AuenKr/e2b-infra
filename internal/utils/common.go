package utils

import (
	commonv1 "e2b/gen/common/v1"

	corev1 "k8s.io/api/core/v1"
)

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
