#!/bin/sh

mkdir -p /tmp/blob/node_1/store

export MANIFEST_DB_FILE=/tmp/blob/node_1/manifest.bolt
export MANIFEST_STORE_DIR=/tmp/blob/node_1/store
export REST_PORT=3001

./blob
