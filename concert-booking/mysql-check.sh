#!/bin/bash

netstat -anpl | grep 3306 >/dev/nul
result=$?

if [[ "$result" != "0" ]]
then
   echo "mysql is unhealthy"
   exit 1
fi

echo "mysql is healthy"
exit 0