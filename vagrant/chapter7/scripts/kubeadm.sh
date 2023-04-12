#!/bin/bash

# configure kubernetes yum repo
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
exclude=kube*
EOF

# disable selinux
setenforce 0
sed -i 's/^SELINUX=.*/SELINUX=permissive/g' /etc/selinux/config

yum install -y  kubernetes-cni-0.6.0 kubelet-1.9.3 kubeadm-1.9.3 kubectl-1.9.3 --disableexcludes=kubernetes

systemctl enable kubelet
systemctl start kubelet
