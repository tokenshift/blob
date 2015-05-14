#!/bin/sh

export BLOB_FILE_SERVICE_PORT=$(( ( RANDOM % 1000 ) + 3000 ))
export BLOB_ADMIN_SERVICE_PORT=$(( ( RANDOM % 1000 ) + 3000 ))

echo  File port: $BLOB_FILE_SERVICE_PORT
echo Admin port: $BLOB_ADMIN_SERVICE_PORT

./blob
