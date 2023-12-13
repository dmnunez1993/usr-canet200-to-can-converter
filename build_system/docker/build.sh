#!/usr/bin/env bash

ARCH=$(uname -m)
USER_UID=$UID
USER_NAME=$USER

cd $ARCH

docker build -f Dockerfile.base -t usr_canet_converter_base .
docker build -f Dockerfile --build-arg USER_UID=$USER_UID --build-arg USER_NAME=$USER_NAME -t usr_canet_converter_builder .
