[[ -z $CODEDEPLOY_DEPLOYMENT_ID ]] && exit 1
[[ -z $CODEDEPLOY_GIT_REPOSITORY ]] && exit 1
[[ -z $CODEDEPLOY_GIT_REFERENCE ]] && exit 1

ls -a /workspace && rm -rf /workspace/* && git clone $CODEDEPLOY_GIT_REPOSITORY \
  --branch $CODEDEPLOY_GIT_REFERENCE \
  --single-branch \
  --depth 1 \
  --no-tags \
  /workspace/$CODEDEPLOY_DEPLOYMENT_ID
