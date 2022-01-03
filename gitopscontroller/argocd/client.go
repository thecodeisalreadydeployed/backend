package argocd

import (
	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
)

func NewArgoCDClient() error {
	client, err := apiclient.NewClient(&apiclient.ClientOptions{
		// TODO(trif0lium): use environment variable
		ServerAddr: "argocd-server.argocd.svc.cluster.local",
		PlainText:  true,
		Insecure:   true,
	})

	if err != nil {
		return err
	}
	closer, sessionClient, err := client.NewSessionClient()

	_ = closer
	_ = sessionClient

	return err
}
