package clients

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func newK8sConfig(in K8sClientParams) (*rest.Config, error) {
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
		return cfg, nil
	}

	// Get from serice account present inside the cluster
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
