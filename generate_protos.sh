#! /usr/bin/env bash

set -e

if ! [[ -x "$(command -v protoc-gen-go)" ]]; then
  echo "Need to install protoc-gen-go"
  exit 1
fi

PROTOC_OPTS="-I./apis/ --go_out=temp --go-grpc_out=temp"

mkdir -p temp/

# shellcheck disable=SC2086
protoc $PROTOC_OPTS ./apis/teams_automator/*.proto

rm -rf pkg/grpc/

mv temp/github.com/kanataidarov/teams_automator/pkg/* pkg

rm -rf temp/