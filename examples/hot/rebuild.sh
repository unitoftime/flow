#!/bin/bash

set -e

PLUGIN_DIR=./plugin/
BUILD_DIR=./plugin/build/
BUILD_PKG_DIR=${BUILD_DIR}/pkg/
BUILD_LIB_DIR=${BUILD_DIR}/lib/

mkdir -p ${BUILD_PKG_DIR}
mkdir -p ${BUILD_LIB_DIR}

rm -f ${BUILD_LIB_DIR}/*.so

VAR=$RANDOM

# Note: I usually leave this uncommented, but I didnt want to accidentally wipe someones computer
#rm -rf ${BUILD_PKG_DIR}/*
mkdir -p ${BUILD_PKG_DIR}/tmp$VAR

cp ${PLUGIN_DIR}/*.go ${BUILD_PKG_DIR}/tmp$VAR

cd ${BUILD_PKG_DIR}/tmp$VAR
sed -i 's/package plugin/package main/g' *.go
cd -

go build -gcflags=all=-d=checkptr=0 -buildmode=plugin -o ${BUILD_LIB_DIR}/tmp$VAR.so ${BUILD_PKG_DIR}/tmp$VAR

echo "Finished Rebuild" $VAR
tree
