#!/usr/bin/env bash

CONTAINER=usr_canet_converter_builder
IMAGE=usr_canet_converter_builder
HOME=~/.usr_canet_converter_build_system_home
CODE_DIR=$(git rev-parse --show-toplevel)
USER_ID=$UID
USER_NAME=$USER
OUTPUT_FOLDER=$1
TARGET=$2

mkdir -p ~/.usr_canet_converter_build_system
mkdir -p ~/.usr_canet_converter_build_system/go

docker run \
  --name $CONTAINER \
  -e USERID=${USER_ID} \
  -e LOCAL_USER_ID=${USER_ID}\
  -v ${HOME}:/home/${USER_NAME}:Z \
  -v ${CODE_DIR}:/usr-canet200-to-can-converter:Z \
  -v ~/.ssh:/home/${USER_NAME}/.ssh \
  -v ${OUTPUT_FOLDER}:/output \
  -w /usr-canet200-to-can-converter \
  --user ${USER_ID}:${USER_ID} \
  -it --rm $IMAGE bash
  # -it --rm $IMAGE bash -c "./build_system/deb_${TARGET}/scripts/build_installer.sh"
