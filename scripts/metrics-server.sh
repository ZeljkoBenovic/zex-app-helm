#!/bin/bash

command -v helm >/dev/null 2>&1 || { echo >&2 "Helm is not installed. Use docs to install: https://helm.sh/docs/intro/install/"; exit 1; }

helm upgrade --install -n metrics --create-namespace -f ms-values.yaml metrics oci://registry-1.docker.io/bitnamicharts/metrics-server
