docker run \
  --pull always \
  -v $(pwd)/testdata/fixture-nest.Dockerfile:/kaniko/deploys-dev/Dockerfile \
  --env CODEDEPLOY_API_URL=http://host.docker.internal:3000 \
  --env CODEDEPLOY_DEPLOYMENT_ID=testdata \
  --env CODEDEPLOY_KANIKO_LOG_VERBOSITY=info \
  --env CODEDEPLOY_KANIKO_CONTEXT=git://github.com/thecodeisalreadydeployed/fixture-nest.git#refs/tags/v2#14bc77fc515e6d66b8d9c15126ee49ca55faf879 \
  ghcr.io/thecodeisalreadydeployed/imagebuilder:latest
