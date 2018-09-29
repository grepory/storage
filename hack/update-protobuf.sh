#!/usr/bin/env bash

ROOT=$(dirname "${BASH_SOURCE}")/..

PACKAGES=(
    github.com/grepory/storage/apis/meta
    github.com/grepory/storage/apis/simple
)

go-to-protobuf \
    --proto-import=./vendor \
    --drop-embedded-fields=github.com/grepory/storage/apis/meta.TypeMeta \
    --packages=$(IFS=, ; echo "${PACKAGES[*]}") \
    --go-header-file=${ROOT}/hack/update-protobuf-boilerplate.txt