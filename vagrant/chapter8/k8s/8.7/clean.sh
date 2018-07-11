# !/bin/bash

kubectl delete -f user.yaml
kubectl delete -f concert.yaml
kubectl delete -f booking.yaml
kubectl delete -f linkerd-internal.yaml
kubectl delete -f linkerd-egress.yaml
kubectl delete -f haproxy.yaml
kubectl delete -f rbac.yaml