#!/bin/bash

# 杀掉旧进程（如果有的话），然后跑个新进程

PROG="$(basename $(dirname $(readlink -f $0)))"

killall -9 ${PROG};
sleep 1;
nohup ./${PROG} 2>&1 &
