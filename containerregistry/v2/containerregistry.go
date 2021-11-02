package containerregistry

type ContainerRegistry interface {
	ImageURL(name string, tag string) string
}

type containerRegistry struct{}
