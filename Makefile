api:
	APP_ENV=DEV go run main.go

dev-image-build:
	DOCKER_BUILDKIT=1 docker build . -t thecodeisalreadydeployed/backend:dev

dev-image-run: dev-image-build
	DATABASE_HOST=localhost && \
	DATABASE_USERNAME=user && \
	DATABASE_PASSWORD=password && \
	DATABASE_NAME=codedeploy && \
	DATABASE_PORT=5432 && \
	docker run thecodeisalreadydeployed/backend:dev

kind:
	sh hack/0-kind-create-cluster.sh
	sh hack/1-kubectl-apply-argocd.sh

port-forward:
	kubectl -n argocd port-forward svc/argocd-server -n argocd 8080:443 &
