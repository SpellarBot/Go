#!/bin/bash

# build项目

PROG="$(basename $(dirname $(readlink -f $0)))"

go build -o ${PROG}