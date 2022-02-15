package types

type Metadata struct {
	Name      string
	Namespace string
	Labels    map[string]string
}

type TLSRedirect struct {
	Enable bool
}

type TLS struct {
	Secret   string
	Redirect TLSRedirect
}

type Upstream struct {
	Name    string
	Service string
	Port    int
}

type RouteAction struct {
	Pass string
}

type Route struct {
	Path   string
	Action RouteAction
}

type Spec struct {
	Host      string
	TLS       TLS
	Upstreams []Upstream
	Routes    []Route
}

type VirtualServer struct {
	APIVersion string
	Kind       string
	Metadata   Metadata
	Spec       Spec
}
