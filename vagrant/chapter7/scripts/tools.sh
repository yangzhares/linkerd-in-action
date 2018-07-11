#/bin/bash

yum install -y wget telnet tree net-tools unzip

# install jq
wget -qO /usr/local/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
chmod +x /usr/local/bin/jq
