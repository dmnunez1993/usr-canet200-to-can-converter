#!/usr/bin/env bash

OUTPUT_FOLDER=$1

cd docker && ./run.sh amd64 $1
