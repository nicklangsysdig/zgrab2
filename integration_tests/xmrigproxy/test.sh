#!/usr/bin/env bash

set -e
MODULE_DIR=$(dirname $0)
ZGRAB_ROOT=$MODULE_DIR/../..
ZGRAB_OUTPUT=$ZGRAB_ROOT/zgrab-output

mkdir -p $ZGRAB_OUTPUT/xmrigproxy

CONTAINER_NAME=zgrab_xmrigproxy

OUTPUT_FILE=$ZGRAB_OUTPUT/xmrigproxy/xmrigproxy.json

echo "xmrigproxy/test: Tests runner for xmrigproxy"
# TODO FIXME: Add any necessary flags or additional tests
CONTAINER_NAME=$CONTAINER_NAME $ZGRAB_ROOT/docker-runner/docker-run.sh xmrigproxy > $OUTPUT_FILE

# Dump the docker logs
echo "xmrigproxy/test: BEGIN docker logs from $CONTAINER_NAME [{("
docker logs --tail all $CONTAINER_NAME
echo ")}] END docker logs from $CONTAINER_NAME"

# TODO: If there are any other relevant log files, dump those to stdout here.
