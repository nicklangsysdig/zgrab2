#!/usr/bin/env bash

set +e

echo "xmrigproxy/cleanup: Tests cleanup for xmrigproxy"

CONTAINER_NAME=zgrab_xmrigproxy

docker stop $CONTAINER_NAME
