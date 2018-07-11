#!/bin/bash
mkdir -p $HOME/.kube

if [[ -f "/etc/kubernetes/admin.conf" ]]
then
    cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    chown $(id -u):$(id -g) $HOME/.kube/config
else
    cp -i /vagrant/admin.conf $HOME/.kube/config
fi
