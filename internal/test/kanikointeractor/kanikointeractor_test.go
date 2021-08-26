package kanikointeractor

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/containerregistry"
	it "github.com/thecodeisalreadydeployed/kanikointeractor"
	"github.com/thecodeisalreadydeployed/util"
)

var kubeconfig = flag.String("kubeconfig", "", "") //nolint

func TestKanikoInteractor_BuildContainerImage(t *testing.T) {
	if os.Getenv("GITHUB_REPOSITORY") != "thecodeisalreadydeployed/backend" {
		t.Skip()
	}

	interactor := it.KanikoInteractor{
		Registry:     containerregistry.LOCAL,
		BuildContext: "https://github.com/thecodeisalreadydeployed/fixture-monorepo.git",
		DeploymentID: util.RandomString(5),
		Destination:  "fixture-nest:dev",
	}

	err := interactor.BuildContainerImage()
	assert.Nil(t, err)
}
