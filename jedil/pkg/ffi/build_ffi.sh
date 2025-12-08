#!/bin/bash
cd /home/zthacker/zach_code/dev-mono/jedil

# Build shared library
go build -buildmode=c-shared -o libjedil.so pkg/ffi/jedil.go

echo "Built libjedil.so"
ls -lh libjedil.so