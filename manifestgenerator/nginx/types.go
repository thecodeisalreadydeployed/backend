package types

type VirtualServer struct {
	APIVersion string
	Kind       string
	Metadata   struct {
		Name      string
		Namespace string
	}
	Spec struct {
		Host string
		TLS  struct {
			Secret   string
			Redirect struct {
				Enable bool
			}
		}
		Upstreams []struct {
			Name    string
			Service string
			Port    int
		}
		Routes []struct {
			Path   string
			Action struct {
				Pass string
			}
		}
	}
}
