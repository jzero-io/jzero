#!/usr/bin/env bash

set -e

Mode=$1

function generateFunc() {
    mkdir -p pkg/client/pb
    # shellcheck disable=SC2043
    for dir in proto; do
      scopeVersion=worktab
      if [ -d "$dir" ]; then
          if [ "$1" = "debug" ]; then
            protoc --proto_path="$dir" --go_out="$(pwd)"/pkg/client/pb --go-httpsdk_out=logtostderr=true,v=1,scopeVersion="$scopeVersion",env_file=./.httpsdk.toml:"$(pwd)"/pkg/client "$dir"/*.proto
          else
            protoc --proto_path="$dir" --go_out="$(pwd)"/pkg/client/pb --go-httpsdk_out=logtostderr=true,scopeVersion="$scopeVersion",env_file=./.httpsdk.toml:"$(pwd)"/pkg/client "$dir"/*.proto
          fi
      fi
    done
    sleep 1
}

generateFunc "$Mode"