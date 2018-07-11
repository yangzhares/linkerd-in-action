#!/bin/bash

kubectl delete -f user.yaml
kubectl delete -f booking.yaml
kubectl delete -f concert.yaml
kubectl delete -f linkerd-tls.yaml
kubectl delete -f rbac.yaml
kubectl delete secret tls
rm -fr certs