#!/bin/sh

set -eux

IS_CONTAINER=${IS_CONTAINER:-false}
CONTAINER_RUNTIME="${CONTAINER_RUNTIME:-podman}"

if [ "${IS_CONTAINER}" != "false" ]; then
  TOP_DIR="${1:-.}"
  export XDG_CACHE_HOME="/tmp/.cache"
  go vet "${TOP_DIR}"/api/... "${TOP_DIR}"/controllers/...
else
  "${CONTAINER_RUNTIME}" run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/metal3-io/hardware-classification-controller:ro,z" \
    --entrypoint sh \
    --workdir /go/src/github.com/metal3-io/hardware-classification-controller \
    registry.hub.docker.com/library/golang:1.16 \
    /go/src/github.com/metal3-io/hardware-classification-controller/hack/govet.sh "${@}"
fi;
