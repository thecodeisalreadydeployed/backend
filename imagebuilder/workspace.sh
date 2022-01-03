[[ -z $CODEDEPLOY_DEPLOYMENT_ID ]] && exit 1
[[ -z $CODEDEPLOY_GIT_REPOSITORY ]] && exit 1
[[ -z $CODEDEPLOY_GIT_REFERENCE ]] && exit 1

rm -rf ./workspace/* && mkdir -p /workspace/$CODEDEPLOY_DEPLOYMENT_ID && cd /workspace/$CODEDEPLOY_DEPLOYMENT_ID && \
  git init && \
  git remote add origin $CODEDEPLOY_GIT_REPOSITORY && \
  git fetch --depth 1 --no-tags --prune --progress origin $CODEDEPLOY_GIT_REFERENCE && \
  git checkout FETCH_HEAD 2>&1 | /deploys-dev/logexporter
