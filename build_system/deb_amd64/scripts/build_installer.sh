#!/usr/bin/env bash

OUTPUT_FOLDER=/output

cd /usr-canet200-to-can-converter/app

echo "Cleaning up previous installation..."
rm -rf $OUTPUT_FOLDER/src

mkdir -p $OUTPUT_FOLDER
mkdir -p $OUTPUT_FOLDER/src/usr_canet200_to_can_converter/usr/local/usr_canet200_to_can_converter

DEBIAN_FOLDER=$OUTPUT_FOLDER/src/usr_canet200_to_can_converter/DEBIAN
ROOT_FOLDER=$OUTPUT_FOLDER/src/usr_canet200_to_can_converter
SRC_FOLDER=$OUTPUT_FOLDER/src/usr_canet200_to_can_converter/usr/local
VAR_FOLDER=$OUTPUT_FOLDER/src/usr_canet200_to_can_converter/var/usr_canet200_to_can_converter
REPOSITORY_ROOT=/usr-canet200-to-can-converter

mkdir -p $VAR_FOLDER

cd /usr-canet200-to-can-converter/app

echo "Installing deps..."
go mod download

echo "Building converter..."
env GOOS=linux GOARCH=amd64 CGO_ENABLED=1\
    go build -o $SRC_FOLDER/usr_canet200_to_can_converter/usr_canet200_to_can_converter cmd/usr_canet200_to_can_converter/usr_canet200_to_can_converter.go

cp -r $REPOSITORY_ROOT/build_system/deb_amd64/metadata/DEBIAN $DEBIAN_FOLDER

echo "Setting up remaining files..."
cp -a $REPOSITORY_ROOT/build_system/deb_amd64/additional/. $ROOT_FOLDER

echo "Creating debian package..."
cd $OUTPUT_FOLDER/src/
fakeroot dpkg-deb --build usr_canet200_to_can_converter $OUTPUT_FOLDER/usr_canet200_to_can_converter.deb
