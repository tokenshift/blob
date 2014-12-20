#!/bin/sh

mkdir -p /tmp/blob/node_1/store

export ADMIN_DB_FILE=/tmp/blob/node_1/admin.bolt
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD_SALT=edc5a5888fd27fa92f526747c1c25ef1e3f3f968a2ac038fcdd5088c3be01814
export ADMIN_PASSWORD_HASH=ab5cded873d9da3f7815b5ff444bb727d074fda17f41241e57f46f54417cbf23
export ADMIN_PORT=4001
export MANIFEST_DB_FILE=/tmp/blob/node_1/manifest.bolt
export MANIFEST_STORE_DIR=/tmp/blob/node_1/store
export REST_PORT=3001

./blob
