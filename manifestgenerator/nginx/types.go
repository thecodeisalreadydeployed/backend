package types

type Metadata struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
}

type TLSRedirect struct {
	Enable bool `json:"enable"`
}

type TLS struct {
	Secret   string      `json:"secret"`
	Redirect TLSRedirect `json:"redirect"`
}

type Upstream struct {
	Name    string `json:"name"`
	Service string `json:"service"`
	Port    int    `json:"port"`
}

type RouteAction struct {
	Pass string `json:"pass"`
}

type Route struct {
	Path   string      `json:"path"`
	Action RouteAction `json:"action"`
}

type Spec struct {
	Host      string     `json:"host"`
	TLS       TLS        `json:"tls"`
	Upstreams []Upstream `json:"upstreams"`
	Routes    []Route    `json:"routes"`
}

type VirtualServer struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}
