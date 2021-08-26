package kanikointeractor

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/containerregistry"
	"github.com/thecodeisalreadydeployed/containerregistry/gcr"
	it "github.com/thecodeisalreadydeployed/kanikointeractor"
	"github.com/thecodeisalreadydeployed/util"
)

var kubeconfig = flag.String("kubeconfig", "", "") //nolint

func TestKanikoInteractor_BuildContainerImage(t *testing.T) {
	registry := gcr.NewGCRGateway("asia.gcr.io", "hu-tao-mains")
	destination, err := registry.RegistryFormat("fixture-monorepo", "dev")
	assert.Nil(t, err)

	interactor := it.KanikoInteractor{
		Registry:     containerregistry.GCR,
		BuildContext: "https://github.com/thecodeisalreadydeployed/fixture-monorepo.git",
		DeploymentID: util.RandomString(5),
		Destination:  destination,
	}

	err = interactor.BuildContainerImage()
	assert.Nil(t, err)
}
