# How to run the integration tests

## Create Kind Cluster
kind create cluster

## Build a1t and a1txapp images
cd onos-a1t
make kind
cd - 

cd onos-a1t/test/utils/xapp
make kind
cd - 

## Local Helm Repos
cd onos-helm-charts && make deps 
cd sdran-helm-charts && make deps

## Public Helm Repos
helm repo add atomix https://charts.atomix.io
helm repo add onosproject https://charts.onosproject.org
helm repo add sdran --username **** --password **** https://sdrancharts.onosproject.org
helm repo update

## Install Atomix Cluster
helm install -n kube-system atomix-controller atomix/atomix-controller --wait
helm install -n kube-system atomix-raft-storage atomix/atomix-raft-storage --wait
helm install -n kube-system onos-operator onos/onos-operator --wait

## Setup a test namespace and bring up CLI and topo
kubectl create namespace test

## Execute the helmit tests
cd onos-a1t
helmit -n test test ./cmd/onos-a1t-test  --secret "sd-ran-username=******" --secret "sd-ran-password=******" --suite a1pm --context ./test/utils/charts/

## Check a1t logs
kubectl -n test logs -f deploy/onos-a1t -c onos-a1t

## Check onos-a1txapp logs
kubectl -n test logs -f deploy/onos-a1txapp -c onos-a1txapp

## As needed, clean the environment
kubectl delete ns test
helm uninstall -n kube-system onos-operator atomix-raft-storage atomix-controller
kind delete cluster kind
