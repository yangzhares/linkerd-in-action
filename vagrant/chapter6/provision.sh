#!/bin/bash

# install docker engine
yum install -y docker-1.13.1
systemctl enable docker
systemctl start docker

# disable selinux
setenforce 0
sed -i 's/^SELINUX=.*/SELINUX=permissive/g' /etc/selinux/config

# install tool sets
yum install -y wget telnet tree net-tools unzip
# install jq
wget -qO /usr/local/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
chmod +x /usr/local/bin/jq

# avoid to erase /etc/resolv.conf when restart networkmanager
grep ^dns=none /etc/NetworkManager/NetworkManager.conf >/dev/null || sed -i '/^plugins/a dns=none' /etc/NetworkManager/NetworkManager.conf
systemctl restart NetworkManager

# install dnsmasq
yum install -y dnsmasq

cat << EOF > /etc/dnsmasq.d/consul
interface=*
addn-hosts=/etc/hosts
bogus-priv
server=/service.consul/127.0.0.1#8600
EOF

systemctl enable dnsmasq && systemctl start dnsmasq
sed -i '/^search/a nameserver 127.0.0.1' /etc/resolv.conf && sed -i 's/\(^search\)/\1 service.consul/g' /etc/resolv.conf # && sed -i '/^search/a nameserver 127.0.0.1' /etc/resolv.conf
chattr +i /etc/resolv.conf
