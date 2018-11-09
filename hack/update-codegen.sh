#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

HACK_DIR=$(dirname "${BASH_SOURCE}")
REPO_ROOT=${HACK_DIR}/..

${REPO_ROOT}/vendor/k8s.io/code-generator/generate-groups.sh \
all \
github.com/jeromefroe/heimdallr/pkg/client \
github.com/jeromefroe/heimdallr/pkg/apis \
heimdallr:v1alpha1 \
--go-header-file ${REPO_ROOT}/hack/boilerplate.go.tmpl \
$@
