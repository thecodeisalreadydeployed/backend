package json

type Deployment struct {
	ApiVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata Metadata `json:"metadata"`
	Spec Spec `json:"spec"`
}

type Metadata struct {
	Name string `json:"name"`
	Labels map[string]string `json:"labels"`
}

type Spec struct {
	Replicas int `json:"replicas"`
	Selector Selector `json:"selector"`
	Template Template `json:"template"`
}

type Selector struct {
	MatchLabels map[string]string `json:"matchLabels"`
}

type Template struct {
	Metadata TemplateMetadata `json:"metadata"`
	Spec TemplateSpec `json:"spec"`
}

type TemplateMetadata struct {
	Labels map[string]string `json:"labels"`
}

type TemplateSpec struct {
	Containers []Container
}

type Container struct {
	Name string
	Image string
	Ports []Port
}

type Port struct {
	ContainerPort []int
}

