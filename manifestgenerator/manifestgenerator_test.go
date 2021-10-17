package manifestgenerator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManifestGenerator(t *testing.T) {

	yamls, err := GenerateManifests()

	if err != nil {
		t.Error(err)
	}

	actualDeploymentYAML := yamls["deploymentYAML"]
	actualServiceYAML := yamls["serviceYAML"]
	actualVirtualServerYAML := yamls["virtualServerYAML"]

	expectedDeploymentYAML := `apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    DeploymentLabelKey: DeploymentLabelValue
  name: DeploymentName
  namespace: DeploymentNamespace
spec:
  selector: null
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        DeploymentLabelKey: DeploymentLabelValue
    spec:
      containers:
      - env:
        - name: PORT
          value: "3000"
        image: DeploymentContainerImage
        imagePullPolicy: IfNotPresent
        name: DeploymentContainerImage
        resources: {}
status: {}
`

	expectedServiceYAML := `apiVersion: apps/v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    ServiceLabelKey: ServiceLabelValue
  name: ServiceName
  namespace: ServiceNamespace
spec: {}
status:
  loadBalancer: {}
`

	expectedVirtualServerYAML := `apiVersion: k8s.nginx.org/v1
kind: VirtualServer
metadata:
  creationTimestamp: null
  labels:
    VirtualServerLabelKey: VirtualServerLabelValue
  name: AppID
  namespace: ProjectID
spec:
  host: AppID.
  http-snippets: ""
  ingressClassName: ""
  policies: null
  routes:
  - action:
      pass: AppID
      proxy: null
      redirect: null
      return: null
    errorPages: null
    location-snippets: ""
    matches: null
    path: /
    policies: null
    route: ""
    splits: null
  server-snippets: ""
  tls:
    redirect: null
    secret: ""
  upstreams:
  - buffer-size: ""
    buffering: null
    buffers: null
    client-max-body-size: ""
    connect-timeout: ""
    fail-timeout: ""
    healthCheck: null
    keepalive: null
    lb-method: ""
    max-conns: null
    max-fails: null
    name: AppID
    next-upstream: ""
    next-upstream-timeout: ""
    next-upstream-tries: 0
    port: 3000
    queue: null
    read-timeout: ""
    send-timeout: ""
    service: AppID
    sessionCookie: null
    slow-start: ""
    subselector: null
    tls:
      enable: false
    use-cluster-ip: false
status:
  message: ""
  reason: ""
  state: ""
`

	assert.Equal(t, expectedDeploymentYAML, actualDeploymentYAML)
	assert.Equal(t, expectedServiceYAML, actualServiceYAML)
	assert.Equal(t, expectedVirtualServerYAML, actualVirtualServerYAML)
}
