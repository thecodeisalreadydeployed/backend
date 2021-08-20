FROM alpine

RUN apk update
RUN apk upgrade
RUN apk add --no-cache openssh-server
RUN ssh-keygen -A
RUN apk add --no-cache git shadow bash

RUN mkdir /__w

COPY init-user.sh /usr/local/bin/init-user
COPY init-repo.sh /usr/local/bin/init-repo
RUN chmod +x /usr/local/bin/init-user /usr/local/bin/init-repo

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
