package kanikogateway

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	gcr "github.com/thecodeisalreadydeployed/containerregistry/gcr"
	gw "github.com/thecodeisalreadydeployed/kanikogateway"
	"github.com/thecodeisalreadydeployed/util"
)

func TestKanikoInteractor_BuildContainerImageGCR(t *testing.T) {
	if os.Getenv("GITHUB_REPOSITORY") == "thecodeisalreadydeployed/backend" && os.Getenv("GITHUB_WORKFLOW") != "kaniko/gcr" {
		t.Skip()
	}

	serviceAccountKey, decodeErr := base64.StdEncoding.DecodeString(os.Getenv("GCP_SERVICE_ACCOUNT_BASE64"))
	if decodeErr != nil {
		fmt.Printf("decodeErr: %v\n", decodeErr)
	}
	assert.Nil(t, decodeErr)

	containerRegistry := gcr.NewGCRGateway("asia.gcr.io", "hu-tao-mains", string(serviceAccountKey))

	appID := util.RandomString(5)
	deploymentID := util.RandomString(5)
	i := gw.NewKanikoGateway(containerRegistry, "https://github.com/thecodeisalreadydeployed/fixture-monorepo.git", appID, deploymentID)

	err := i.BuildContainerImage()
	assert.Nil(t, err)
}