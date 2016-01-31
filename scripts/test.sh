#!/bin/bash
set -e -x -u

ROOT_DIR_PATH=$(cd $(dirname $0)/.. && pwd)
cd $ROOT_DIR_PATH

ginkgo ./client ./commands ./controller ./shell ./aws/templates
