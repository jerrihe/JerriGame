#!/bin/bash

OUT_PATH=../../src/protocol/cs

if [ ! -d $OUT_PATH ]; then
    mkdir -p $OUT_PATH
fi

protoc --go_out=$OUT_PATH ./*.proto