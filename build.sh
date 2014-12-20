#!/bin/sh

mdtangle *.go.md **/*.go.md **/**/*.go.md && \
go build && \
go install github.com/tokenshift/blob/admin/bhash
