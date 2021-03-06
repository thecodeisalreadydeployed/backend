docker run \
  -it \
  --pull always \
  -v $(pwd)/testdata/fixture-nest.Dockerfile:/kaniko/deploys-dev/Dockerfile \
  --env CODEDEPLOY_API_URL=http://host.docker.internal:3000 \
  --env CODEDEPLOY_DEPLOYMENT_ID=testdata \
  --entrypoint /bin/sh \
  ghcr.io/thecodeisalreadydeployed/imagebuilder:latest
