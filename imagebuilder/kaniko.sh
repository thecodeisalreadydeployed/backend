/kaniko/executor \
  --no-push \
  --log-format=json \
  --verbosity=debug \
  --dockerfile=/kaniko/deploys-dev/Dockerfile \
  --context=git://github.com/thecodeisalreadydeployed/fixture-nest.git#refs/tags/v2 2>&1 | tee /kaniko/deploys-dev/kaniko.log | /kaniko/deploys-dev/logexporter && echo "DONE" && cat /kaniko/deploys-dev/kaniko.log
