package argocd

import "net/http"

func NewArgoCDClient() error {
	apiURL := "http://argocd-server.argocd.svc.cluster.local/api/v1/application?name=" + "codedeploy" + "&refresh=true"
	req, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	_ = req
	return nil
}
