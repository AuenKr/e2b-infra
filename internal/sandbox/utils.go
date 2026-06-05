package sandbox

import (
	commonv1 "e2b/gen/common/v1"

	corev1 "k8s.io/api/core/v1"
)

func podPhaseToSandboxStateAdapter(status corev1.PodPhase) commonv1.SandboxStatus {
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
