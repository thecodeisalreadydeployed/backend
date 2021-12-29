FROM gcr.io/kaniko-project/executor:v1.7.0-debug
COPY ./logexporter-linux-amd64 /kaniko/deploys-dev/logexporter
COPY ./imagebuilder/kaniko.sh /kaniko/deploys-dev/kaniko.sh
ENTRYPOINT /bin/sh /kaniko/deploys-dev/kaniko.sh
