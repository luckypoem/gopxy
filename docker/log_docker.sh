#!/bin/bash

id=`docker ps -a | grep "gopxy" | awk -F" " '{print $1}'`

docker logs -f -t --tail 10 $id