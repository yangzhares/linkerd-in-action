#!/bin/bash

# 由于linkerd作为sidecar，与user服务部署在一起,
# 需注册linkerd到Consul, 所有linkerd均以相同名字linkerd注册，
# 使用端口及标签分辩linkerd服务于哪种服务，演示环境中,
# 约定4141=>user服务,4142=>booking服务,4143=>concert服务
SERVICE_NAME=linkerd
SERVICE_TAG=user
DOCKER_NAME=$SERVICE_TAG-$(cat /dev/urandom | head -n 10 | md5sum | head -c 10)
IP=$(ip addr show | grep eth1 | grep inet | awk '{print $2}' | cut -d'/' -f1)

# user服务配置信息, 可通过环境变量或JSON文件进行配置
DBNAME=demo
DBUSER=test
DBPASSWORD=pass
DBENDPOINT=mysql.service.consul:3306
# user服务访问booking和concert服务时使用如下地址，
# [tag].[service-name].service.consul
BOOKING_SERVICE_ADDR=booking.linkerd.service.consul:4142
CONCERT_SERVICE_ADDR=concert.linkerd.service.consul:4143

# 启动user服务，并将linkerd通过Registrator注册到Consul
docker run -d \
    -p 4141:4141 \
    -p 62000:8180 \
    --dns $IP \
    --name $DOCKER_NAME \
    --env DBNAME=$DBNAME \
    --env DBUSER=$DBUSER \
    --env DBPASSWORD=$DBPASSWORD \
    --env DBENDPOINT=$DBENDPOINT \
    --env BOOKING_SERVICE_ADDR=$BOOKING_SERVICE_ADDR \
    --env CONCERT_SERVICE_ADDR=$CONCERT_SERVICE_ADDR \
    --env SERVICE_4141_NAME=$SERVICE_NAME \
    --env SERVICE_TAGS=$SERVICE_TAG \
    --env SERVICE_8180_IGNORE=true \
    --volume $PWD/linkerd-user.yml:/linkerd.yml \
    zhanyang/user:1.1
