package sandbox

import (
	"context"
	"fmt"
	"strings"

	commonv1 "e2b/gen/common/v1"
	sandboxv1 "e2b/gen/sandbox/v1"
	"e2b/pkg/config"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// StartSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) StartSandbox(ctx context.Context, req *sandboxv1.StartSandboxRequest) (*sandboxv1.StartSandboxResponse, error) {
	// Start the pod with given image
	k8core := s.K8sClient.CoreV1()

	// Create Pod
	cpuLimit, memoryLimit := config.CPU_MAX_DEFAULT, config.MEMORY_MAX_DEFAULT

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
			corev1.ResourceCPU:    resource.MustParse(config.CPU_MIN_DEFAULT),
			corev1.ResourceMemory: resource.MustParse(config.MEMORY_MIN_DEFAULT),
		},
	}

	cmds := strings.Fields(req.GetCmd())
	args := strings.Fields(req.GetArgs())

	envs := make([]corev1.EnvVar, 0, len(req.GetEnv()))
	for name, value := range req.GetEnv() {
		envs = append(envs, corev1.EnvVar{
			Name:  name,
			Value: value,
		})
	}

	container := corev1.Container{
		Name:      req.Id,
		Image:     req.Image,
		Command:   cmds,
		Args:      args,
		Ports:     []corev1.ContainerPort{},
		Env:       envs,
		Resources: resources,
		// LivenessProbe:   &corev1.Probe{},
		// ReadinessProbe:  &corev1.Probe{},
		// StartupProbe:    &corev1.Probe{},
		ImagePullPolicy: "IfNotPresent",
		SecurityContext: &corev1.SecurityContext{},
	}

	podConfg := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.GetId(),
			Namespace: s.Config.K8sNamespace,
			Labels: map[string]string{
				"id": req.GetId(),
			},
		},
		Spec: corev1.PodSpec{
			Containers:      []corev1.Container{container},
			SecurityContext: &corev1.PodSecurityContext{},
			Hostname:        req.GetId(),
		},
	}

	podInfo, err := k8core.Pods(s.Config.K8sNamespace).Create(ctx, &podConfg, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	// Create service
	serviceSpec := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Id,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"id": req.Id,
			},
			Ports: []corev1.ServicePort{
				{
					Name:     config.INITIAL_PORT_NAME,
					Protocol: K8sProtocolAdapter(config.INITIAL_PORT_PROTOCOL),
					Port:     config.INITIAL_PORT,
				},
			},
		},
	}
	_, err = k8core.Services(s.Config.K8sNamespace).Create(ctx, &serviceSpec, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	// Create HTTP Route
	gatewayNamespace := gatewayapiv1.Namespace(s.Config.K8sGatewayNamespace)
	httpSpec := gatewayapiv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Id,
		},
		Spec: gatewayapiv1.HTTPRouteSpec{
			CommonRouteSpec: gatewayapiv1.CommonRouteSpec{
				ParentRefs: []gatewayapiv1.ParentReference{
					{
						Namespace: &gatewayNamespace,
						Name:      gatewayapiv1.ObjectName(s.Config.K8sGateway),
					},
				},
			},
		},
	}
	_, err = s.GatewayClient.GatewayV1().HTTPRoutes(s.Config.K8sNamespace).Create(ctx, &httpSpec, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return &sandboxv1.StartSandboxResponse{
		Sandbox: &commonv1.SandboxInfo{
			Id:     podInfo.Name,
			Status: PodPhaseToSandboxStateAdapter(podInfo.Status.Phase),
		},
	}, nil
}

// StopSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) StopSandbox(ctx context.Context, req *sandboxv1.StopSandboxRequest) (*sandboxv1.StopSandboxResponse, error) {
	k8core := s.K8sClient.CoreV1()

	// Delete pod
	if err := k8core.Pods(s.Config.K8sNamespace).Delete(ctx, req.GetId(), metav1.DeleteOptions{}); err != nil {
		return nil, err
	}
	// Delete service
	if err := k8core.Services(s.Config.K8sNamespace).Delete(ctx, req.GetId(), metav1.DeleteOptions{}); err != nil {
		return nil, err
	}
	// Delete HTTP Route
	if err := s.GatewayClient.GatewayV1().HTTPRoutes(s.Config.K8sNamespace).Delete(ctx, req.GetId(), metav1.DeleteOptions{}); err != nil {
		return nil, err
	}
	return &sandboxv1.StopSandboxResponse{
		Sandbox: &commonv1.SandboxInfo{
			Id:     req.Id,
			Status: commonv1.SandboxStatus_SANDBOX_STATUS_STOPPED,
		},
	}, nil
}

// GetSandbox implements [sandboxv1connect.SandboxServiceHandler].
func (s *SandboxServer) GetSandbox(ctx context.Context, req *sandboxv1.GetSandboxRequest) (*sandboxv1.GetSandboxResponse, error) {
	k8core := s.K8sClient.CoreV1()

	podInfo, err := k8core.Pods(s.Config.K8sNamespace).Get(ctx, req.GetId(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var cpuMilli, memoryBytes int64

	// If it fails, metric will be nil. But our api will return the pod info.
	podMetrics, err := s.MetricsClient.MetricsV1beta1().PodMetricses(s.Config.K8sNamespace).Get(ctx, req.GetId(), metav1.GetOptions{})
	if err != nil {
		s.Logger.Error("failed to get metrics", zap.Error(err))
	} else {
		for _, container := range podMetrics.Containers {
			cpuMilli += container.Usage.Cpu().MilliValue()
			memoryBytes += container.Usage.Memory().Value()
		}
	}

	//

	return &sandboxv1.GetSandboxResponse{
		Sandbox: &commonv1.SandboxInfo{
			Id:     podInfo.Name,
			Status: PodPhaseToSandboxStateAdapter(podInfo.Status.Phase),
		},
		Resource: &commonv1.Specification{
			Cpu:    uint32(cpuMilli),
			Memory: uint32(memoryBytes / (1024 * 1024)), // In MB
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
	sandboxs := make([]*commonv1.SandboxInfo, len(podsInfo.Items))
	for i, podInfo := range podsInfo.Items {
		sandboxs[i] = &commonv1.SandboxInfo{
			Id:     podInfo.Name,
			Status: PodPhaseToSandboxStateAdapter(podInfo.Status.Phase),
		}
	}

	return &sandboxv1.ListSandboxResponse{
		Sandboxs: sandboxs,
	}, nil
}
