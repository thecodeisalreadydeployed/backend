#!/bin/sh

if [ "$#" != 2 ]
then
  exit 1
fi

USERNAME=$1
REPO=$2.git
W=/__w

PATH=$W/$USERNAME/$REPO

ls -a /usr/bin
git init --bare $PATH
