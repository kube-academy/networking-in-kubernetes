#/bin/bash
echo ------------------- Front end deployment -----------------------
kubectl get all -n front-end
echo
echo ------------------- Back end deployment -----------------------
kubectl get all -n back-end
