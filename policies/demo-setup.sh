#!/usr/bin/env bash

# Starting mysql db outside of cluster
docker-compose up -d

kind create cluster --config kind-calico.yaml
kubectl apply -f calico.yaml

docker build -t demoapp:latest .
kind load docker-image demoapp:latest

kubectl create ns back-end
kubectl create configmap dbhost -n back-end --from-literal dbhost=$(hostname -I | cut -f1 -d' ')
kubectl apply -f app.yaml

sleep 10s
kubectl port-forward service/web -n front-end --address 0.0.0.0 9000:80