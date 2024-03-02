#!/bin/bash

command -v kind >/dev/null 2>&1 || { echo >&2 "Kind cli is not installed. Use docs to install: https://kind.sigs.k8s.io/docs/user/quick-start/#installing-with-a-package-manager"; exit 1; }

kind create cluster --config kind_cluster_deploy.yaml
