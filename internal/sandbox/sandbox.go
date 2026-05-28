package sandbox

import (
	"context"
	"fmt"

	sandboxv1 "e2b/gen/sandbox/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CPU_MIN_DEFAULT    = "50m"
	MEMORY_MIN_DEFAULT = "10Mi"
)

// StartSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) StartSandbox(ctx context.Context, req *sandboxv1.StartSandboxRequest) (*sandboxv1.StartSandboxResponse, error) {
	// Start the pod with given image
	k8core := s.K8sClient.CoreV1()

	cpuLimit, memoryLimit := CPU_MIN_DEFAULT, MEMORY_MIN_DEFAULT

	if req.Requirement != nil {
		cpuLimit = fmt.Sprintf("%vm", req.Requirement.Cpu)
		memoryLimit = fmt.Sprintf("%vMi", req.Requirement.Memory)
	}
	resources := corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuLimit),
			corev1.ResourceMemory: resource.MustParse(memoryLimit),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(CPU_MIN_DEFAULT),
			corev1.ResourceMemory: resource.MustParse(MEMORY_MIN_DEFAULT),
		},
	}
	container := corev1.Container{
		Name:      req.Id,
		Image:     req.Image,
		Command:   []string{},
		Args:      []string{},
		Ports:     []corev1.ContainerPort{},
		Env:       []corev1.EnvVar{},
		Resources: resources,
		// LivenessProbe:   &corev1.Probe{},
		// ReadinessProbe:  &corev1.Probe{},
		// StartupProbe:    &corev1.Probe{},
		ImagePullPolicy: "IfNotPresent",
		SecurityContext: &corev1.SecurityContext{},
	}

	podConfg := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.GetId(),
			Namespace:   s.Config.K8sNamespace,
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		},
		Spec: corev1.PodSpec{
			Volumes:             []corev1.Volume{},
			InitContainers:      []corev1.Container{},
			Containers:          []corev1.Container{container},
			EphemeralContainers: []corev1.EphemeralContainer{},
			SecurityContext:     &corev1.PodSecurityContext{},
			Hostname:            req.GetId(),
			Affinity:            &corev1.Affinity{},
			Tolerations:         []corev1.Toleration{},
		},
	}
	opts := metav1.CreateOptions{}

	podInfo, err := k8core.Pods(s.Config.K8sNamespace).Create(ctx, &podConfg, opts)
	if err != nil {
		return nil, err
	}

	return &sandboxv1.StartSandboxResponse{
		Sandbox: &sandboxv1.SandboxInfo{
			Id:     podInfo.Name,
			Status: convertPodPhaseToSandboxState(podInfo.Status.Phase),
			Url:    podInfo.Status.PodIP,
		},
	}, nil
}

// StopSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) StopSandbox(ctx context.Context, req *sandboxv1.StopSandboxRequest) (*sandboxv1.StopSandboxResponse, error) {
	// Stop the pod
	k8core := s.K8sClient.CoreV1()

	if err := k8core.Pods(s.Config.K8sNamespace).Delete(ctx, req.GetId(), metav1.DeleteOptions{}); err != nil {
		return nil, err
	}
	return &sandboxv1.StopSandboxResponse{
		Sandbox: &sandboxv1.SandboxInfo{
			Id:     req.Id,
			Status: sandboxv1.SandboxStatus_SANDBOX_STATUS_STOPPED,
			Url:    "",
		},
	}, nil
}

// GetSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) GetSandbox(ctx context.Context, req *sandboxv1.GetSandboxRequest) (*sandboxv1.GetSandboxResponse, error) {
	k8core := s.K8sClient.CoreV1()

	podInfo, err := k8core.Pods(s.Config.K8sNamespace).Get(ctx, req.Id, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &sandboxv1.GetSandboxResponse{
		Sandbox: &sandboxv1.SandboxInfo{
			Id:     podInfo.Name,
			Status: convertPodPhaseToSandboxState(podInfo.Status.Phase),
			Url:    podInfo.Status.PodIP,
		},
		Resource: &sandboxv1.Specification{
			Cpu:    uint32(podInfo.Spec.Overhead.Cpu().AsApproximateFloat64()),
			Memory: uint32(podInfo.Spec.Overhead.Memory().AsApproximateFloat64()),
			// Storage: podInfo.Sp,
		},
	}, nil
}

// ListSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) ListSandbox(ctx context.Context, req *sandboxv1.ListSandboxRequest) (*sandboxv1.ListSandboxResponse, error) {
	// Get pod details
	k8core := s.K8sClient.CoreV1()

	podsInfo, err := k8core.Pods(s.Config.K8sNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	sandboxs := make([]*sandboxv1.SandboxInfo, len(podsInfo.Items))
	for i, podInfo := range podsInfo.Items {
		sandboxs[i] = &sandboxv1.SandboxInfo{
			Id:     podInfo.Name,
			Status: convertPodPhaseToSandboxState(podInfo.Status.Phase),
			Url:    podInfo.Status.PodIP,
		}
	}

	return &sandboxv1.ListSandboxResponse{
		Sandboxs: sandboxs,
	}, nil
}
