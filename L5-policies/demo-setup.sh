#!/usr/bin/env bash

# Starting mysql db outside of cluster
docker-compose up -d

kind create cluster --config kind-calico.yaml
kubectl apply -f calico.yaml

docker build -t demoapp:latest .
kind load docker-image demoapp:latest

kubectl create ns back-end
kubectl create configmap dbhost -n back-end --from-literal dbhost=$(hostname -I | cut -f1 -d' ')

kubectl create ns front-end
kubectl create configmap dbhost -n front-end --from-literal dbhost=$(hostname -I | cut -f1 -d' ')

kubectl apply -f app.yaml

echo 'alias k="kubectl"' > ~/.bashrc
