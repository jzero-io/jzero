#!/bin/bash

set -e

WorkPath=$(pwd)

FmtType=$1

# export tools bin to PATH
PATH="${WorkPath}"/tools:"${PATH}"

if [ ! -f "$WorkPath"/tools/goimports-reviser-darwin ]; then
  {
    mkdir -p tools
    curl -o "${WorkPath}"/tools/goimports-reviser-darwin https://oss.jaronnie.com/goimports-reviser-darwin
    chmod +x "${WorkPath}"/tools/goimports-reviser-darwin
  }
fi

if [ "$FmtType" = "all" ]; then
  echo "start to goimports go file, please wait..."
  for i in $(find . ! -path "./_example/*" ! -path "./scopesourcedir/*" -name "*.go"  | grep -v ".pb.go" | grep -v ".pb.gw.go") ; do
    echo "goimports for $i"
    "$WorkPath"/tools/goimports-reviser-darwin -rm-unused -set-alias -format -file-path "$i"
  done
  echo "successfully goimports"
else
  for i in $(git diff --cached --name-only --diff-filter=ACM -- '*.go' | grep -v ".pb.go" | grep -v ".pb.gw.go") ; do
      "$WorkPath"/tools/goimports-reviser-darwin -rm-unused -set-alias -format -file-path "$i"
      git add "$i"
  done
fi