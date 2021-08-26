package kanikointeractor

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/containerregistry"
	gcr "github.com/thecodeisalreadydeployed/containerregistry/gcr"
	it "github.com/thecodeisalreadydeployed/kanikointeractor"
	"github.com/thecodeisalreadydeployed/util"
)

var kubeconfig = flag.String("kubeconfig", "", "") //nolint

func TestKanikoInteractor_BuildContainerImage(t *testing.T) {
	interactor := it.KanikoInteractor{
		Registry:     containerregistry.LOCAL,
		BuildContext: "https://github.com/thecodeisalreadydeployed/fixture-monorepo.git",
		DeploymentID: util.RandomString(5),
		Destination:  "fixture-nest:dev",
	}

	err := interactor.BuildContainerImage()
	assert.Nil(t, err)
}

func TestKanikoInteractor_BuildContainerImageGCR(t *testing.T) {
	gateway := gcr.NewGCRGateway("asia.gcr.io", "hu-tao-mains")
	destination, err := gateway.RegistryFormat("fixture-monorepo", "dev")
	assert.Nil(t, err)

	interactor := it.KanikoInteractor{
		Registry:     containerregistry.LOCAL,
		BuildContext: "https://github.com/thecodeisalreadydeployed/fixture-monorepo.git",
		DeploymentID: util.RandomString(5),
		Destination:  destination,
	}

	err = interactor.BuildContainerImage()
	assert.Nil(t, err)
}
