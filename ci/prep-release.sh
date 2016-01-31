#!/bin/bash
set -e -u -x

export OUT_BINARIES=$PWD/bin
export OUT_NOTES=$PWD/notes

version=$(cat version/number)
echo "v${version}" > $OUT_NOTES/name

mkdir -p go/src/github.com/rosenhouse/
cp -a proctor-source go/src/github.com/rosenhouse/proctor
cd go

export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin
export GO15VENDOREXPERIMENT=1

cd src/github.com/rosenhouse/proctor

for os in linux darwin windows; do
  GOOS=$os go build -o $OUT_BINARIES/proctor-$os
done

git rev-parse HEAD > $OUT_NOTES/commitish
