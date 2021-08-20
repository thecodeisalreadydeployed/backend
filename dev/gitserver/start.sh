#!/bin/sh

if [ "$(ls -A /__w/keys/)" ]; then
  cd /home/git
  cat /__w/keys/*.pub > .ssh/authorized_keys
  chown -R codedeploy:codedeploy .ssh
  chmod 700 .ssh
  chmod -R 600 .ssh/*
fi

if [ "$(ls -A /__w/repos/)" ]; then
  cd /__w/repos
  chown -R codedeploy:codedeploy .
  chmod -R ug+rwX .
  find . -type d -exec chmod g+s '{}' +
fi

/usr/sbin/sshd -D
