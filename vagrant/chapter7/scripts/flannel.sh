#/bin/bash

IFNAME=$1

curl -s https://raw.githubusercontent.com/coreos/flannel/v0.10.0/Documentation/kube-flannel.yml > kube-flannel.yml
sed -ri 's/^(\s*)(\"Type\"\s*: \"vxlan\"\s*$)/\1"Type": "host-gw"/' kube-flannel.yml
sed -i "/- --kube-subnet-mgr/a \ \ \ \ \ \ \ \ - --iface=${IFNAME}" kube-flannel.yml

kubectl apply -f kube-flannel.yml
