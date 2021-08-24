package kanikointeractor

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/containerregistry"
	it "github.com/thecodeisalreadydeployed/kanikointeractor"
	"github.com/thecodeisalreadydeployed/util"
)

var kubeconfig = flag.String("kubeconfig", "", "")

func TestKanikoInteractor_BuildContainerImage(t *testing.T) {
	interactor := it.KanikoInteractor{
		Registry:     containerregistry.LOCAL,
		BuildContext: "https://github.com/thecodeisalreadydeployed/fixture-nest.git",
		DeploymentID: util.RandomString(5),
	}

	err := interactor.BuildContainerImage()
	assert.Nil(t, err)
}