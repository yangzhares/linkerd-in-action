#!/bin/bash

# 调用user服务API创建用户信息并返回
user_id=$(curl \
    -s \
    -X POST \
    -d '{"ID": "tom","Name": "Tom Gao","Age": 23}' \
    http://10.244.1.38:8180/users | jq -r '.id')

# 调用concert服务API创建演唱会信息并返回concert ID用作预定使用
concert_id=$(curl \
    -s \
    -X POST \
    -d '{"concert_name": "The best of Andy Lau 2018","singer": "Andy Lau","start_date": "2018-03-27 20:30:00","end_date": "2018-04-07 23:00:00","location": "Shanghai","street": "Jiangwan Stadium"}' \
    http://10.244.1.39:8182/concerts | jq -r '.id')

# 调用booking服务API预定演唱会
curl \
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
    http://10.244.2.32:8181/bookings >/dev/null


curl \
    -s \
    -X POST \
    -d '{"concert_name": "The best of Andy Lau 2019","singer": "Jack Zhang","start_date": "2018-04-27 20:30:00","end_date": "2018-05-07 23:00:00","location": "Shanghai","street": "Jiangwan Stadium"}' \
    http://10.244.2.27:8182/concerts