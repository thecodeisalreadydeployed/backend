FROM gcr.io/kaniko-project/executor:v1.6.0-debug
COPY ./logexporter-linux-amd64 /kaniko/deploys-dev/logexporter
COPY ./imagebuilder/kaniko.sh /kaniko/deploys-dev/kaniko.sh
RUN ["/busybox/mkdir", "-p", "/bin"]
RUN ["/busybox/ln", "-s", "/busybox/sh", "/bin/sh"]
ENTRYPOINT /bin/sh /kaniko/deploys-dev/kaniko.sh
