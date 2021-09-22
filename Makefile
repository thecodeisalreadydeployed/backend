api:
	APP_ENV=DEV go run main.go

dev-image:
	docker-compose -f docker-compose.yaml -f docker-compose.dev.yaml up --force-recreate -- backend

kind:
	sh hack/0-kind-create-cluster.sh
	sh hack/1-kubectl-apply-argocd.sh

port-forward:
	kubectl -n argocd port-forward svc/argocd-server -n argocd 8080:443 &
