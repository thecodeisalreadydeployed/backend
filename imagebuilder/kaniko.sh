/kaniko/executor \
  --no-push \
  --log-format=json \
  --verbosity=debug \
  --dockerfile=/kaniko/deploys-dev/Dockerfile \
  --context=dir:///workspace 2>&1 | tee /kaniko/deploys-dev/kaniko.log | /kaniko/deploys-dev/logexporter && echo "DONE" && cat /kaniko/deploys-dev/kaniko.log
