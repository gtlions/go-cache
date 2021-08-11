#!/bin/bash
rm -rf ../cache/*;protoc --go_out=../cache --go_opt=paths=source_relative --go-grpc_out=../cache --go-grpc_opt=paths=source_relative cache.proto