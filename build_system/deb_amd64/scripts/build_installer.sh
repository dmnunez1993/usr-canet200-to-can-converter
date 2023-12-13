#!/usr/bin/env bash

OUTPUT_FOLDER=/output

cd /usr-canet200-to-can-converter/app

echo "Cleaning up previous layout..."
rm -rf $OUTPUT_FOLDER/src
rm -rf $OUTPUT_FOLDER/admin

mkdir -p $OUTPUT_FOLDER
mkdir -p $OUTPUT_FOLDER/src/usr_canet200_to_can_converter/usr/local/usr_canet200_to_can_converter

DEBIAN_FOLDER=$OUTPUT_FOLDER/src/usr_canet200_to_can_converter/DEBIAN
ROOT_FOLDER=$OUTPUT_FOLDER/src/usr_canet200_to_can_converter
SRC_FOLDER=$OUTPUT_FOLDER/src/usr_canet200_to_can_converter/usr/local
VAR_FOLDER=$OUTPUT_FOLDER/src/usr_canet200_to_can_converter/var/usr_canet200_to_can_converter
REPOSITORY_ROOT=/usr-canet200-to-can-converter

mkdir -p $VAR_FOLDER

# Build Golang App
cd /usr-canet200-to-can-converter/app

echo "Installing deps..."
go mod download

echo "Building converter..."
env GOOS=linux GOARCH=amd64 CGO_ENABLED=1\
    go build -o $SRC_FOLDER/usr_canet200_to_can_converter/usr_canet200_to_can_converter cmd/usr_canet200_to_can_converter/usr_canet200_to_can_converter.go

# Build NodeJS admin
cd /usr-canet200-to-can-converter
# Exclude files that might be present if developing on the same machine
rsync -av --progress admin $OUTPUT_FOLDER --exclude admin/node_modules/ --exclude admin/build/
cd $OUTPUT_FOLDER/admin
yarn install
yarn build
mkdir -p $SRC_FOLDER/usr_canet200_to_can_converter/admin
cp -R build/. $SRC_FOLDER/usr_canet200_to_can_converter/admin/

# Set up remaining files required to build the package
cp -r $REPOSITORY_ROOT/build_system/deb_amd64/metadata/DEBIAN $DEBIAN_FOLDER

echo "Setting up remaining files..."
cp -a $REPOSITORY_ROOT/build_system/deb_amd64/additional/. $ROOT_FOLDER

echo "Creating debian package..."
cd $OUTPUT_FOLDER/src/
fakeroot dpkg-deb --build usr_canet200_to_can_converter $OUTPUT_FOLDER/usr_canet200_to_can_converter.deb
