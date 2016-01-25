#!/bin/bash
set -e -u -x

mkdir -p go/src/github.com/rosenhouse/
cp -a bosh-proctor go/src/github.com/rosenhouse/bosh-proctor
cd go

export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin
export GO15VENDOREXPERIMENT=1

cd src/github.com/rosenhouse/bosh-proctor

go install ./vendor/github.com/onsi/ginkgo/ginkgo

./scripts/test.sh
