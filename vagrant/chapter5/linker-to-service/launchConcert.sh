#!/bin/bash

# 由于linkerd作为sidecar，与concert服务部署在一起,
# 需注册linkerd到Consul, 所有linkerd均以相同名字linkerd注册，
# 使用端口及标签分辩linkerd服务于哪种服务，演示环境中,
# 约定4141=>user服务,4142=>booking服务,4143=>concert服务
SERVICE_NAME=linkerd
SERVICE_TAG=concert
DOCKER_NAME=$SERVICE_TAG-$(cat /dev/urandom | head -n 10 | md5sum | head -c 10)
IP=$(ip addr show | grep eth1 | grep inet | awk '{print $2}' | cut -d'/' -f1)

# concert, 可通过环境变量或JSON文件进行配置
DBNAME=demo
DBUSER=root
DBPASSWORD=pass
DBENDPOINT=mysql.service.consul:3306

# 启动concert服务，并将linkerd通过Registrator注册到Consul
docker run -d \
    -p 4143:4141 \
    -p 62002:8182 \
    --dns $IP \
    --name $DOCKER_NAME \
    --env SERVICE_4141_NAME=$SERVICE_NAME \
    --env SERVICE_TAGS=$SERVICE_TAG \
    --env SERVICE_8182_IGNORE=true \
    --env DBNAME=$DBNAME \
    --env DBUSER=$DBUSER \
    --env DBPASSWORD=$DBPASSWORD \
    --env DBENDPOINT=$DBENDPOINT \
    --volume $PWD/linkerd-concert.yml:/linkerd.yml \
    zhanyang/concert:1.1
