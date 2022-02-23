set -o pipefail

[[ -z $CODEDEPLOY_DEPLOYMENT_ID ]] && exit 1
[[ -z $CODEDEPLOY_KANIKO_LOG_VERBOSITY ]] && CODEDEPLOY_KANIKO_LOG_VERBOSITY="debug"
[[ -z $CODEDEPLOY_KANIKO_LOG_FILE ]] && CODEDEPLOY_KANIKO_LOG_FILE="/kaniko/deploys-dev/kaniko.log"
[[ -z $CODEDEPLOY_KANIKO_CONTEXT ]] && CODEDEPLOY_KANIKO_CONTEXT="dir:///workspace"
[[ -z $CODEDEPLOY_KANIKO_DOCKERFILE ]] && CODEDEPLOY_KANIKO_DOCKERFILE="/kaniko/deploys-dev/Dockerfile"

if [[ -z $CODEDEPLOY_KANIKO_DESTINATION ]]; then
  /kaniko/executor \
    --no-push \
    --force \
    --log-format=json \
    --verbosity=$CODEDEPLOY_KANIKO_LOG_VERBOSITY \
    --dockerfile=$CODEDEPLOY_KANIKO_DOCKERFILE \
    --context=$CODEDEPLOY_KANIKO_CONTEXT 2>&1 | tee $CODEDEPLOY_KANIKO_LOG_FILE | /kaniko/deploys-dev/logexporter && echo "DONE" && cat $CODEDEPLOY_KANIKO_LOG_FILE
else
  /kaniko/executor \
    --destination=$CODEDEPLOY_KANIKO_DESTINATION \
    --force \
    --cache=true \
    --use-new-run \
    --log-format=json \
    --verbosity=$CODEDEPLOY_KANIKO_LOG_VERBOSITY \
    --dockerfile=$CODEDEPLOY_KANIKO_DOCKERFILE \
    --context=$CODEDEPLOY_KANIKO_CONTEXT 2>&1 | tee $CODEDEPLOY_KANIKO_LOG_FILE | /kaniko/deploys-dev/logexporter && echo "DONE" && cat $CODEDEPLOY_KANIKO_LOG_FILE
fi
