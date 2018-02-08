#!/usr/bin/env bash

builddir=`mktemp -d`
IMAGE="10.199.192.16/machine_learning/mlcloud:v1.0.1"

echo "build context: $builddir"

cp build.sh Dockerfile $builddir

cp -r ../../src $builddir

cp -r ../../vendor $builddir

cd $builddir && docker build . -t $IMAGE

echo "built image $IMAGE"
echo "remove temp build dir"
rm -rf $builddir
