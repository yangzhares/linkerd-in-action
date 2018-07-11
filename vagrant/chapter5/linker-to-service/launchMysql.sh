#!/bin/bash

MYSQL_ROOT_PASSWORD=pass
MYSQL_DATABASE=demo
MYSQL_USER=test
MYSQL_PASSWORD=pass

SERVICE_NAME=mysql
IP=$(ip addr show | grep eth1 | grep inet | awk '{print $2}' | cut -d'/' -f1)

docker run -d \
    -p 3306:3306 \
    --name mysql \
    --network host \
    --env MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
    --env MYSQL_DATABASE=$MYSQL_DATABASE \
    --env MYSQL_USER=$MYSQL_USER \
    --env MYSQL_PASSWORD=$MYSQL_PASSWORD \
    --env IP=$IP \
    --env SERVICE_NAME=$SERVICE_NAME \
    zhanyang/mysql:8.0
