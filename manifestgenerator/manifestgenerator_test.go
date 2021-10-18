package manifestgenerator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateDeployment(t *testing.T) {
	options := &GenerateDeploymentOptions{
		Name:           "DeploymentName",
		Namespace:      "DeploymentNamespace",
		Labels:         map[string]string{"DeploymentLabelKey": "DeploymentLabelValue"},
		ContainerImage: "DeploymentContainerImage",
	}

	expected := `apiVersion: apps/v1
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
	actual, err := GenerateDeploymentYAML(options)

	assert.Nil(t, err)
	assert.YAMLEq(t, expected, actual)
}

func TestGenerateService(t *testing.T) {
	options := &GenerateServiceOptions{
		Name:      "ServiceName",
		Namespace: "ServiceNamespace",
		Labels:    map[string]string{"ServiceLabelKey": "ServiceLabelValue"},
	}

	expected := `apiVersion: apps/v1
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
	actual, err := GenerateServiceYAML(options)
	assert.Nil(t, err)
	assert.YAMLEq(t, expected, actual)
}
