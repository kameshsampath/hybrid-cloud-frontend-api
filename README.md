# Hybrid Cloud Front End

The front end API that will be used with Gloo Hybrid Cloud Demo

## Pre-requisites

* [Docker Desktop](https://docs.docker.com/desktop/)
* [pipx](https://pypa.github.io/pipx)
* [kubectl](https://kubernetes.io/docs/tasks/tools)
* [httpie](https://httpie.io)
* Kubernetes Cluster e.g [KinD](https://kind.sigs.k8s.io)
* [Gloo Edge](https://docs.solo.io/gloo-edge/latest/getting_started/

## Deploy API

```shell
kubectl apply -k k8s/app
```