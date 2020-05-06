#!/bin/sh

set -eux

IS_CONTAINER=${IS_CONTAINER:-false}
ARTIFACTS=${ARTIFACTS:-/tmp}
CONTAINER_RUNTIME="${CONTAINER_RUNTIME:-podman}"

if [ "${IS_CONTAINER}" != "false" ]; then
  eval "$(go env)"
  cd "${GOPATH}"/src/github.com/metal3-io/hardware-classification-controller
  export XDG_CACHE_HOME="/tmp/.cache"
  go test -v ./pkg/... ./cmd/... -coverprofile "${ARTIFACTS}"/cover.out
else
  "${CONTAINER_RUNTIME}" run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/metal3-io/hardware-classification-controller:ro,z" \
    --entrypoint sh \
    --workdir /go/src/github.com/metal3-io/hardware-classification-controller \
    registry.hub.docker.com/library/golang:1.12 \
    /go/src/github.com/metal3-io/hardware-classification-controller/hack/unit.sh "${@}"
fi;
