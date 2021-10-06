.PHONY: dev
dev:
	docker compose -f docker-compose.yaml up -d
	APP_ENV=DEV go run main.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-ci
	sh hack/lint.sh

.PHONY: e2e
e2e:
	go test -v ./test

.PHONY: dev-image
dev-image:
	DOCKER_BUILDKIT=1 docker compose -f docker-compose.yaml -f docker-compose.dev.yaml up --build -- backend

.PHONY: dev-image-detach
dev-image-detach:
	DOCKER_BUILDKIT=1 docker compose -f docker-compose.yaml -f docker-compose.dev.yaml up --build --detach -- backend

.PHONY: kind
kind:
	sh hack/0-kind-create-cluster.sh
	sh hack/1-kubectl-apply-argocd.sh

.PHONY: port-forward
port-forward:
	@kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d && echo
	kubectl -n argocd port-forward svc/argocd-server -n argocd 8080:443
