api:
	docker compose -f docker-compose.yaml up -d
	APP_ENV=DEV go run main.go

e2e:
	go test -v ./test

dev-image:
	DOCKER_BUILDKIT=1 docker compose -f docker-compose.yaml -f docker-compose.dev.yaml up --build -- backend

dev-image-detach:
	DOCKER_BUILDKIT=1 docker compose -f docker-compose.yaml -f docker-compose.dev.yaml up --build --detach -- backend

kind:
	sh hack/0-kind-create-cluster.sh
	sh hack/1-kubectl-apply-argocd.sh

port-forward:
	@kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d && echo
	kubectl -n argocd port-forward svc/argocd-server -n argocd 8080:443
