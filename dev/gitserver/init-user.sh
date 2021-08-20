#!/bin/sh

if [ "$#" != 2 ]
then
  exit 1
fi

USERNAME=$1
PASSWORD=$2
W=/__w

mkdir -p $W/$USERNAME
