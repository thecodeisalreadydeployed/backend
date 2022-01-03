FROM alpine/git:v2.30.2
COPY ./logexporter-linux-amd64 /deploys-dev/logexporter
COPY ./imagebuilder/workspace.sh /deploys-dev/workspace.sh
ENTRYPOINT /bin/sh /deploys-dev/workspace.sh
