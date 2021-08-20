#!/bin/sh

if [ "$#" != 3 ]
then
  exit 1
fi

USERNAME=$1
PASSWORD=$2
REPO=$3.git
W=/__w

mkdir -p $W/$USERNAME

PATH=$W/$USERNAME/$REPO

ls -a /usr/bin
git init --bare $PATH
