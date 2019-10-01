#!/bin/bash

set -e

#==========args here==============
BUILD_VERSION=v0.0.1
BIND_PORT=8080

#=================================
mkdir -p conf
cp config.json ./conf

if [  ! -f ./conf/ca.key.pem ] || [ ! -f ./conf/ca.pem ]; then
  echo "not found cert file, create it..."
  DIR=`pwd`
  cd ../certs/
  ./openssl-gen.sh
  mv *.pem $DIR/conf
  cd $DIR
  echo "create cert file finish, recommend to add it to your browser."
else
  echo "found cert file, skip create..."
fi

echo "begin clean docker instance..."
./clean_docker.sh

echo "begin build docker image..."
docker build --build-arg VERSION=$BUILD_VERSION -t "gopxy/gopxy" .

echo "begin run container..."
docker run -it -d --restart=always --name "gopxy" -p $BIND_PORT:8080 -v `pwd`/conf:/etc/gopxy "gopxy/gopxy"

echo "all finish..."


