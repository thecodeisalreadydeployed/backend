package main

import (
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/thecodeisalreadydeployed/kubernetesinteractor"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func absolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	if strings.HasPrefix(path, string(0x7E)) {
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		return filepath.Join(home, strings.TrimPrefix(path, string(0x7E)))
	}

	return path
}

func main() {
	var kubeconfig string
	pflag.StringVar(&kubeconfig, "kubeconfig", "", "")
	pflag.Parse()
	kubeconfig = absolutePath(kubeconfig)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	kubernetesinteractor.ListDeployments(clientset)

	return
}
