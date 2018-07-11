# !/bin/bash

kubectl delete -f user.yaml
kubectl delete -f concert.yaml
kubectl delete -f booking.yaml
kubectl delete -f linkerd-ingress-controller.yaml
kubectl delete -f linkerd-ingress-controller-tls.yaml
kubectl delete -f linkerd.yaml
kubectl delete -f rbac.yaml
kubectl delete -f user-ingress.yaml
kubectl delete secret ingress-cert
rm -fr certs