#!/bin/bash

set -e
IFNAME=$1
IP="$(ip -4 a s $IFNAME | grep "inet" | head -1 |awk '{print $2}' | cut -d'/' -f1)"
sed -e "s/^.*${HOSTNAME}.*/${IP} ${HOSTNAME} ${HOSTNAME}/" -i /etc/hosts
