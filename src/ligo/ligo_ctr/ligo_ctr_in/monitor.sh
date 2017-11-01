#!/bin/bash

# 进程保活

PROG="$(basename $(dirname $(readlink -f $0)))"
DIR="$(dirname $(readlink -f $0))"
DINGTALK_API="https://oapi.dingtalk.com/robot/send?access_token=b7b9d6858ff9ab06bc9adb8113b8c0e55817908b0fd07ca1b2cba6965f50fd82"
DINGTALK_CTYPE="Content-Type: application/json"
MACHINE=`hostname -s`

date
ulimit -c unlimited
ps -ef|grep ${PROG}|grep -v grep|grep -v monitor
DEAMON_NUM=`ps -ef|grep ${PROG}|grep -v grep|grep -v monitor|wc -l`
echo $DEAMON_NUM

if [ $DEAMON_NUM -lt 1 ] ; then
  MSG="Try to Restart ${PROG} in ${MACHINE}"
  echo ${MSG}
  curl -s "${DINGTALK_API}" -H "${DINGTALK_CTYPE}" -d "{\"msgtype\": \"text\", \"text\": { \"content\": \"${MSG}\" } }"
  cd ${DIR}
  nohup ./${PROG} 2>&1 &
  sleep 1
  ps -ef|grep ${PROG}|grep -v grep|grep -v monitor
fi

