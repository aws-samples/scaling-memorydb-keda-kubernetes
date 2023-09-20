#!/bin/bash

kubectl apply -f https://github.com/kedacore/keda/releases/download/v2.10.0/keda- 2.10.0.yaml
kubectl get crd
kubectl get deployment -n keda
kubectl logs -f $(kubectl get pod -l=app=keda-operator -o jsonpath='{.items[0].metadata.name}' -n keda) -n keda
