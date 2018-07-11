#!/bin/bash

# 调用user服务API创建用户信息并返回
ip=$(kubectl get svc user -o json | jq -r '.spec.clusterIP')
port=$(kubectl get svc user -o json | jq '.spec.ports[]|select(.name == "http")|.port')
user_id=$(/bin/curl \
    -s \
    -X POST \
    -d '{"ID": "tom","Name": "Tom Gao","Age": 23}' \
    http://$ip:$port/users | jq -r '.id')

# 调用concert服务API创建演唱会信息并返回concert ID用作预定使用
ip=$(kubectl get svc concert -o json | jq -r '.spec.clusterIP')
port=$(kubectl get svc concert -o json | jq '.spec.ports[]|select(.name == "http")|.port')
concert_id=$(/bin/curl \
    -s \
    -X POST \
    -d '{"concert_name": "The best of Andy Lau 2018","singer": "Andy Lau","start_date": "2018-03-27 20:30:00","end_date": "2018-04-07 23:00:00","location": "Shanghai","street": "Jiangwan Stadium"}' \
    http://$ip:$port/concerts | jq -r '.id')

# 调用booking服务API预定演唱会
ip=$(kubectl get svc booking -o json | jq -r '.spec.clusterIP')
port=$(kubectl get svc booking -o json | jq '.spec.ports[]|select(.name == "http")|.port')
/bin/curl \
    -s \
    -X POST \
    -d @<(cat <<EOF
{
    "user_id": "${user_id}",
    "date": "2018-04-02 20:30:00",
    "concert_id": "${concert_id}"
}
EOF
    ) \
    http://$ip:$port/bookings >/dev/null