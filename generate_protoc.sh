#!/bin/bash
go get github.com/chyngyz-sydykov/book-rating-protos@latest

PROTO_PATH=$(go list -m -f '{{.Dir}}' github.com/chyngyz-sydykov/book-rating-protos)

OUTPUT_DIR=./proto

protoc --proto_path="$PROTO_PATH" \
       --go_out="$OUTPUT_DIR" \
       --go-grpc_out="$OUTPUT_DIR" \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       "$PROTO_PATH/rating/rating.proto"

echo gRPC files are generated in "$OUTPUT_DIR"