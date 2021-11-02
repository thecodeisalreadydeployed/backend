package containerregistry

type ContainerRegistry interface {
	ImageName(name string, tag string) string
}

type containerRegistry struct{}
