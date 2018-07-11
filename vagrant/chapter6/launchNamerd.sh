#!/bin/bash

SERVICE_NAME=namerd
IP=$(ip addr show | grep eth1 | grep inet | awk '{print $2}' | cut -d'/' -f1)

docker run -d \
    -p 4321:4321 \
		-p 9991:9991 \
		-p 4180:4180 \
    --name namerd \
    --network host \
    --env IP=$IP \
    --env SERVICE_4321_NAME=$SERVICE_NAME \
    --env SERVICE_9991_IGNORE=true \
    --env SERVICE_4180_IGNORE=true \
		--volume $PWD/namerd.yml:/namerd.yml \
    buoyantio/namerd:1.3.6 /namerd.yml

# wait Namerd service ready
sleep 10

./namerctl dtab \
	create demo $PWD/demo.dtab \
	--base-url http://namerd.service.consul:4180
