FROM alpine

RUN apk update
RUN apk upgrade
RUN apk add --no-cache openssh git
RUN ssh-keygen -A
WORKDIR /__w/
RUN mkdir /__w/keys \
  && adduser -D -s /usr/bin/git-shell git \
  && echo codedeploy:password | chpasswd \
  && mkdir /home/git/.ssh
COPY sshd_config /etc/ssh/sshd_config
COPY start.sh start.sh

EXPOSE 22
CMD ["sh", "start.sh"]
