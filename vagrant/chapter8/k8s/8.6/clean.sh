# !/bin/bash

kubectl delete -f user.yaml
kubectl delete -f concert.yaml
kubectl delete -f booking.yaml
kubectl delete -f linkerd-edge.yaml
kubectl delete -f linkerd.yaml
kubectl delete -f rbac.yaml
kubectl delete -f haproxy.yaml
kubectl delete -f haproxy-sni.yaml
kubectl delete secret haproxy-certs
rm -fr certs