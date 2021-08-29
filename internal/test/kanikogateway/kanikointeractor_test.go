package kanikointeractor

import (
	"flag"
	"os"
	"testing"
)

var kubeconfig = flag.String("kubeconfig", "", "") //nolint

func TestKanikoInteractor_BuildContainerImage(t *testing.T) {
	if os.Getenv("GITHUB_REPOSITORY") != "thecodeisalreadydeployed/backend" {
		t.Skip()
	}
}