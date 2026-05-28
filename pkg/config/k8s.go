package config

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
	gatewayclientset "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

type K8sClientParams struct {
	fx.In
	Config Config
	Logger *zap.Logger
}

func NewK8sRESTConfig(in K8sClientParams) (*rest.Config, error) {
	return newK8sConfig(in)
}

func NewK8sClusterClient(in K8sClientParams) (*kubernetes.Clientset, error) {
	k8sConfig, err := newK8sConfig(in)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	in.Logger.Info("Successfully created cluster client")
	return client, nil
}

func NewK8sMetricClient(in K8sClientParams) (*metricsclientset.Clientset, error) {
	k8sConfig, err := newK8sConfig(in)
	if err != nil {
		return nil, err
	}

	client, err := metricsclientset.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	in.Logger.Info("Successfully created cluster metrics client")
	return client, nil
}

func NewK8sGatewayClient(in K8sClientParams) (*gatewayclientset.Clientset, error) {
	k8sConfig, err := newK8sConfig(in)
	if err != nil {
		return nil, err
	}

	client, err := gatewayclientset.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	in.Logger.Info("Successfully created cluster gateway client")
	return client, nil
}
