package main

import (
	"encoding/json"
	"net/http"
	"os"

	types "github.com/thecodeisalreadydeployed/metadataserver/types"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	})

	http.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	})

	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	})

	http.ListenAndServe(":5000", nil)
}
