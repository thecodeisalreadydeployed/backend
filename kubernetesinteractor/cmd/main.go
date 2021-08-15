package main

import "github.com/thecodeisalreadydeployed/kubernetesinteractor"

func main() {
	it := kubernetesinteractor.NewKubernetesInteractor()
	it.ListDeployments()
	return
}
