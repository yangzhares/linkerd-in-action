#!/bin/bash

cat <<EOF >  /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables = 1
EOF

sysctl --system
