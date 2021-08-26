package kanikointeractor

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	gcr "github.com/thecodeisalreadydeployed/containerregistry/gcr"
	it "github.com/thecodeisalreadydeployed/kanikointeractor"
	"github.com/thecodeisalreadydeployed/util"
)

var serviceAccountKey = flag.String("gcp-service-account", "", "") //nolint

func TestKanikoInteractor_BuildContainerImageGCR(t *testing.T) {
	if os.Getenv("GITHUB_REPOSITORY") == "thecodeisalreadydeployed/backend" && os.Getenv("GITHUB_WORKFLOW") != "kaniko/gcr" {
		t.Skip()
	}

	gateway := gcr.NewGCRGateway("asia.gcr.io", "hu-tao-mains", *serviceAccountKey)
	destination, err := gateway.RegistryFormat("fixture-monorepo", "dev")
	assert.Nil(t, err)

	fmt.Printf("destination: %v\n", destination)

	interactor := it.KanikoInteractor{
		Registry:     gateway,
		BuildContext: "https://github.com/thecodeisalreadydeployed/fixture-monorepo.git",
		DeploymentID: util.RandomString(5),
		Destination:  destination,
	}

	err = interactor.BuildContainerImage()
	assert.Nil(t, err)
}
