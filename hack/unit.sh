#!/bin/sh

set -eux

IS_CONTAINER=${IS_CONTAINER:-false}
CONTAINER_RUNTIME="${CONTAINER_RUNTIME:-docker}"

if [ "${IS_CONTAINER}" != "false" ]; then
  export XDG_CACHE_HOME=/tmp/.cache
  mkdir /tmp/unit
  cp -r . /tmp/unit
  cd /tmp/unit
  make test
else
  "${CONTAINER_RUNTIME}" run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/metal3-io/hardware-classification-controller:ro,z" \
    --entrypoint sh \
    --workdir /go/src/github.com/metal3-io/hardware-classification-controller \
    quay.io/metal3-io/capm3-unit:master \
    /go/src/github.com/metal3-io/hardware-classification-controller/hack/unit.sh "${@}"
fi;
