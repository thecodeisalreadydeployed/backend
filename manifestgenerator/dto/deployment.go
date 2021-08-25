package dto

type Deployment struct {
	APIVersion    string
	Name          string
	Replicas      int
	Labels        map[string]string
	ContainerSpec ContainerSpec
}

type ContainerSpec struct {
	Name  string
	Image string
	Port  int
}
