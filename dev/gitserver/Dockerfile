FROM alpine:3.7

RUN apk add openssh-server
RUN ssh-keygen -A
RUN apk add git
RUN apk add shadow

VOLUME /__w
WORKDIR /__w

COPY adduser.sh /usr/local/bin/adduser
COPY gitinit.sh /usr/local/bin/gitinit
RUN chmod +x /usr/local/bin/adduser /usr/local/bin/gitinit

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
