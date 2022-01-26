FROM gcr.io/kaniko-project/executor:a7425d1fd0442b58dc24698285102176365a28d9-debug
COPY ./logexporter-linux-amd64 /kaniko/deploys-dev/logexporter
COPY ./imagebuilder/kaniko.sh /kaniko/deploys-dev/kaniko.sh
ENTRYPOINT /bin/sh /kaniko/deploys-dev/kaniko.sh
