#!/bin/bash

did=`docker ps -a | grep "gopxy" | awk -F" " '{print $1}'`
for id in $did
do
  echo "found docker id:"$id" clean it..."
  docker stop $id
  docker rm $id
done