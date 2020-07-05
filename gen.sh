#!/bin/sh

#CLIENT_OUTDIR=frontend-proto/template
SERVER_OUTPUT_DIR=src/proto

mkdir -p ${SERVER_OUTPUT_DIR}

protoc --proto_path=protocol -Iprotocol -Iprotocol/vendor/protobuf/src ./protocol/*.proto \
    --go_out=plugins=grpc:${SERVER_OUTPUT_DIR}

