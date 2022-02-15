package main

import (
	"net/http"
	"os"
)

type metadataHandler struct{}

func (handler *metadataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func main() {
	nodeName := os.Getenv("NODE_NAME")
	podName := os.Getenv("POD_NAME")
	podNamespace := os.Getenv("POD_NAMESPACE")
	projectID := os.Getenv("PROJECT_ID")
	appID := os.Getenv("APP_ID")
	deploymentID := os.Getenv("DEPLOYMENT_ID")

	mux := http.NewServeMux()
	mux.Handle("/", &metadataHandler{})
	http.ListenAndServe(":5000", mux)
}
