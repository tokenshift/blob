#!/bin/sh

export BLOB_FILE_SERVICE_PORT=3000
export BLOB_FILE_STORE_DB="db_1.bdb"
export BLOB_FILE_STORE_DIR="$(pwd)/db_1/"
export BLOB_ADMIN_SERVICE_PORT=4000

mkdir -p $BLOB_FILE_STORE_DIR

echo "     File Port: $BLOB_FILE_SERVICE_PORT"
echo "    Admin Port: $BLOB_ADMIN_SERVICE_PORT"
echo " File Store DB: $BLOB_FILE_STORE_DB"
echo "File Store Dir: $BLOB_FILE_STORE_DIR"

./blob