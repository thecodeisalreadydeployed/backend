package containerregistry

type gcrGateway struct {
}

func NewGCRGateway() ContainerRegistry {
	return gcrGateway{}
}

func (gateway gcrGateway) ImageURL(name string, tag string) string {
	return ""
}
