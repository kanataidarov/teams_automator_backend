#! /usr/bin/env bash

set -e

if ! [[ -x "$(command -v protoc-gen-go)" ]]; then
  echo "Need to install protoc-gen-go"
  exit 1
fi

PROTOC_OPTS="-I./apis/ --go_out=temp --go-grpc_out=temp"

mkdir -p temp/

# shellcheck disable=SC2086
protoc $PROTOC_OPTS ./apis/interview_automator/*.proto

rm -rf pkg/interview_automator

mv temp/github.com/kanataidarov/interview_automator/pkg/* pkg

rm -rf temp/