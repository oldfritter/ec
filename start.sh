#!/bin/sh

# nohup ./cmd/workers >> logs/workers.log 2>&1 &
nohup ./cmd/api >> logs/api.log 2>&1 &
