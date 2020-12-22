#!/bin/sh
go build -ldflags '-w -s' -o ./cmd/workers services/worker/workers.go
go build -ldflags '-w -s' -o ./cmd/api services/api/api.go
