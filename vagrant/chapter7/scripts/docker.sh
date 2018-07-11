#!/bin/bash

yum install -y docker-1.13.1
systemctl enable docker && systemctl start docker
