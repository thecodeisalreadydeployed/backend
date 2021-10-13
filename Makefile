.PHONY: dcp
dcp:
	docker compose -f docker-compose.yaml up -d

.PHONY: dev
dev: dcp
	APP_ENV=DEV go run main.go

.PHONY: prod
prod: dcp
	DATABASE_HOST=localhost \
	DATABASE_USERNAME=user \
	DATABASE_PASSWORD=password \
	DATABASE_NAME=codedeploy \
	DATABASE_PORT=5432 \
	APP_ENV=PROD \
	go run main.go

.PHONY: docs
docs:
	cd design-docs/ && yarn
	cd design-docs/ && yarn start

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -v ./test

.PHONY: lint-ci
lint-ci:
	sh hack/lint.sh

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
