FROM alpine:3.7

RUN apk add openssh-server
RUN ssh-keygen -A
RUN apk add git
RUN apk add shadow

VOLUME /__w
WORKDIR /__w
