#!/usr/bin/env bash
NODE=$(kubectl get node -o json | jq -r .items[0].status.addresses[0].address)

kubectl exec -it -n front-end $(kubectl get pod -n front-end -l run=web -o name) -- curl -s --connect-timeout 2 localhost:9000/health

echo
kubectl exec -it -n back-end $(kubectl get pod -n back-end -l run=dbsvc -o name) -- curl -s --connect-timeout 2 localhost:9000/health

echo
kubectl exec -it -n back-end $(kubectl get pod -n back-end -l run=extsvc -o name) -- curl -s --connect-timeout 2 localhost:9000/health

echo
echo Outside to web
curl -s --connect-timeout 2 ${NODE}:30000

echo
echo Outside to dbsvc
curl -s --connect-timeout 2 ${NODE}:30001

echo
echo Outside to extsvc
curl -s --connect-timeout 2 ${NODE}:30002

