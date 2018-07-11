# !/bin/bash

kubectl delete -f user-1.0.yaml
kubectl delete -f user-2.0.yaml
kubectl delete -f concert.yaml
kubectl delete -f booking.yaml
kubectl delete -f linkerd.yaml
kubectl delete -f linkerd-namerd.yaml
kubectl delete -f namerd.yaml
kubectl delete -f rbac.yaml
kubectl delete -f rbac-namerd.yaml