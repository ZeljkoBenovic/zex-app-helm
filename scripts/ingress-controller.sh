#!/bin/bash

command -v helm >/dev/null 2>&1 || { echo >&2 "Helm is not installed. Use docs to install: https://helm.sh/docs/intro/install/"; exit 1; }

helm install -n ingress-nginx --create-namespace --repo https://kubernetes.github.io/ingress-nginx ingress-nginx --values ingress_nginx.yaml --version 4.10.0 ingress-nginx
