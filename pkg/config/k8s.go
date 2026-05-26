package config

import (
	"os"
	"path/filepath"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClientParams struct {
	fx.In
	Config Config
	Logger *zap.Logger
}

func NewK8sClusterClient(in K8sClientParams) (*kubernetes.Clientset, error) {
	var k8sConfig *rest.Config
	if in.Config.Mode == "dev" {
		// Location path for kubeconfig
		// The default location for the kubeconfig file is in the user's home directory.
		var kubeconfig string
		if home := os.Getenv("HOME"); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
		if in.Config.K8sConfigPath != "" {
			kubeconfig = in.Config.K8sConfigPath
		}

		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		k8sConfig = cfg
	} else {
		// Get from serice account present inside the cluster
		cfg, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		k8sConfig = cfg
	}

	client, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	in.Logger.Info("Successfully connected to the cluster")
	return client, nil
}
