# How to run the integration tests

## Create Kind Cluster
```bash
kind create cluster
```

## Build a1t and a1txapp images
```bash
cd onos-a1t
make kind
cd - 
```

```bash
cd onos-a1t/test/utils/xapp
make kind
cd - 
```

## Public Helm Repos
```bash
helm repo add atomix https://charts.atomix.io
helm repo add onosproject https://charts.onosproject.org
helm repo add sdran https://sdrancharts.onosproject.org
helm repo update
```

## Install Atomix Cluster
```bash
helm install atomix-controller atomix/atomix-controller -n kube-system --wait --version 0.6.8
helm install atomix-raft-storage atomix/atomix-raft-storage -n kube-system --wait --version 0.1.15
helm install onos-operator onos/onos-operator -n kube-system --wait --version 0.4.14 
```

## Setup a test namespace and bring up CLI and topo
```bash
kubectl create namespace test
```

## Execute the helmit tests

```bash
cd onos-a1t
helmit -n test test ./cmd/onos-a1t-test --suite a1pm --context ./test/utils/charts/
```

## Check a1t logs
```bash
kubectl -n test logs -f deploy/onos-a1t -c onos-a1t
```

## Check onos-a1txapp logs
```bash
kubectl -n test logs -f deploy/onos-a1txapp -c onos-a1txapp
```

## As needed, clean the environment
```bash
kubectl delete ns test
helm uninstall -n kube-system onos-operator atomix-raft-storage atomix-controller
kind delete cluster
```
