package main

import (
	"encoding/json"
	"net/http"
	"os"

	types "github.com/thecodeisalreadydeployed/metadataserver/types"
)

type metadataHandler struct{}

func (handler *metadataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nodeName := os.Getenv("NODE_NAME")
	podName := os.Getenv("POD_NAME")
	podNamespace := os.Getenv("POD_NAMESPACE")
	projectID := os.Getenv("PROJECT_ID")
	appID := os.Getenv("APP_ID")
	deploymentID := os.Getenv("DEPLOYMENT_ID")

	metadataBytes, _ := json.Marshal(types.Metadata{
		NodeName:     nodeName,
		PodName:      podName,
		PodNamespace: podNamespace,
		ProjectID:    projectID,
		AppID:        appID,
		DeploymentID: deploymentID,
	})

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(metadataBytes)

}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &metadataHandler{})
	http.ListenAndServe(":5000", mux)
}
