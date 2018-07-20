#!/bin/bash

# Copyright Greg Haskins All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

set -e

# place the Go build cache directory into the default build tree if it exists
if [ -d "${GOPATH}/src/github.com/mcc-github/blockchain/.build" ]; then
    export GOCACHE="${GOPATH}/src/github.com/mcc-github/blockchain/.build/go-cache"
fi

blockchain_dir="$(cd "$(dirname "$0")/.." && pwd)"
source_dirs=$(go list -f '{{.Dir}}' ./... | sed s,"${blockchain_dir}".,,g | cut -f 1 -d / | sort -u)

echo "Checking with gofmt"
OUTPUT="$(gofmt -l -s ${source_dirs})"
if [ -n "$OUTPUT" ]; then
    echo "The following files contain gofmt errors"
    echo "$OUTPUT"
    echo "The gofmt command 'gofmt -l -s -w' must be run for these files"
    exit 1
fi

echo "Checking with goimports"
OUTPUT="$(goimports -l ${source_dirs} | grep -Ev '(^|/)testdata/' || true)"
if [ -n "$OUTPUT" ]; then
    echo "The following files contain goimports errors"
    echo $OUTPUT
    echo "The goimports command 'goimports -l -w' must be run for these files"
    exit 1
fi

echo "Checking with go vet"
PRINTFUNCS="Print,Printf,Info,Infof,Warning,Warningf,Error,Errorf,Critical,Criticalf,Sprint,Sprintf,Log,Logf,Panic,Panicf,Fatal,Fatalf,Notice,Noticef,Wrap,Wrapf,WithMessage"
OUTPUT="$(go vet -printfuncs $PRINTFUNCS ./...)"
if [ -n "$OUTPUT" ]; then
    echo "The following files contain go vet errors"
    echo $OUTPUT
    exit 1
fi
