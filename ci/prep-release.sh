#!/bin/bash
set -e -u -x

export OUT_BINARIES=$PWD/binaries
export OUT_NOTES=$PWD/release-notes

version=$(cat version/number)
echo "v${version}" > $OUT_NOTES/name

mkdir -p go/src/github.com/rosenhouse/
cp -a proctor-source go/src/github.com/rosenhouse/proctor
cd go

export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin
export GO15VENDOREXPERIMENT=1

cd src/github.com/rosenhouse/proctor

for GOOS in linux darwin windows; do
  go build -o $OUT_BINARIES/proctor-$GOOS &
done

wait

git config --global user.email "$GIT_USER_EMAIL"
git config --global user.name "$GIT_USER_NAME"

git rev-parse HEAD > $OUT_NOTES/commitish
