#!/bin/bash

set -ex

OUTPUT=hack/tools/bin/controller-gen

# Check for a vendor directory if any downstream forks use that dependency
# tracking method
if [ -d "vendor" ]
then
    go build -o "${OUTPUT}" ./vendor/sigs.k8s.io/controller-tools/cmd/controller-gen
    exit 0;
fi


CMDPATH="$GOPATH/pkg/mod/sigs.k8s.io/controller-tools@v0.4.0/cmd/controller-gen"

if [ ! -f "$CMDPATH" ]
then
    go mod download
fi

go build -o "${OUTPUT}" "$CMDPATH"
