package dto

type Service struct {
	ApiVersion string
	Name       string
	Labels     map[string]string
	PortSpec   PortSpec
}

type PortSpec struct {
	Protocol   string
	Port       int
	TargetPort int
}
