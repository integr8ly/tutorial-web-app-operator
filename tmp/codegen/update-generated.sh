#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

vendor/k8s.io/code-generator/generate-groups.sh \
deepcopy \
github.com/integr8ly/tutorial-web-app-operator/pkg/generated \
github.com/integr8ly/tutorial-web-app-operator/pkg/apis \
integreatly:v1alpha1 \
--go-header-file "./tmp/codegen/boilerplate.go.txt"
