FROM alpine

RUN apk update
RUN apk upgrade
RUN apk add --no-cache openssh git
RUN ssh-keygen -A
WORKDIR /__w/
RUN mkdir /__w/keys \
  && adduser -D -s /usr/bin/git-shell codedeploy \
  && echo codedeploy:password | chpasswd \
  && mkdir /home/codedeploy/.ssh
COPY sshd_config /etc/ssh/sshd_config
COPY git-shell-commands /home/codedeploy/git-shell-commands
COPY authorized-keys.sh /usr/local/bin/authorized-keys.sh
COPY start.sh start.sh

EXPOSE 22
CMD ["sh", "start.sh"]
