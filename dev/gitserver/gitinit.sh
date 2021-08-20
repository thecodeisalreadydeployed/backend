#!/bin/sh

if [ "$#" != 2 ]
then
  exit 1
fi

USERNAME=$1
REPO=$2.git
W=/__w

PATH=$W/$USERNAME/$REPO

git init --bare $PATH

chown $USERNAME.$USERNAME -R $PATH
chmod -R 700 $PATH
