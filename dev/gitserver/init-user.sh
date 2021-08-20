#!/bin/sh

if [ "$#" != 2 ]
then
  exit 1
fi

USERNAME=$1
PASSWORD=$2
W=/__w

adduser -D -h "/home/$USERNAME" -s '/usr/bin/git-shell' $USERNAME
echo "$USERNAME:$PASSWORD" | chpasswd

mkdir -p $W/$USERNAME

chown $USERNAME.$USERNAME -R $W/$USERNAME
chmod 700 -R $W/$USERNAME
