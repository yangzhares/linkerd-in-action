#!/bin/bash

function rand(){
    start=$1
    end=$(($2-$start+1))
    num=$(cat /dev/urandom | head -n 10 | cksum | awk -F ' ' '{print $1}')
    echo $(($num%$end+$start))
}

# 随机生成服务端口以及注册服务需要的信息
PORT=$(rand 30000 65535)
SERVICE_NAME=booking
DOCKER_NAME=$SERVICE_NAME-$(cat /dev/urandom | head -n 10 | md5sum | head -c 10)
IP=$(ip addr show | grep eth1 | grep inet | awk '{print $2}' | cut -d'/' -f1)

# booking服务配置信息, 可通过环境变量或JSON文件进行配置
DBNAME=demo
DBUSER=root
DBPASSWORD=pass
DBENDPOINT=mysql.service.consul:3306
CONCERT_SERVICE_ADDR=concert.service.consul
SERVICE_ENDPOINT="$IP:$PORT"

# 启动booking服务，并通过Registrator注册到Consul
# 注入http_proxy=localhost:4140使得所有HTTP请求转发到linkerd
docker run -d \
    -p $PORT:$PORT \
    --network host \
    --dns $IP \
    --name $DOCKER_NAME \
    --env http_proxy=localhost:4140 \
    --env IP=$IP \
    --env SERVICE_NAME=$SERVICE_NAME \
    --env DBNAME=$DBNAME \
    --env DBUSER=$DBUSER \
    --env DBPASSWORD=$DBPASSWORD \
    --env DBENDPOINT=$DBENDPOINT \
    --env CONCERT_SERVICE_ADDR=$CONCERT_SERVICE_ADDR \
    --env SERVICE_ENDPOINT=$SERVICE_ENDPOINT \
    zhanyang/booking:1.0