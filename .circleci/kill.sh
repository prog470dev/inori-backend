#!/bin/sh

#var1=`ps -ef | grep ino.out | grep -v grep | wc -l`

#echo $var1

#if [ $var1 > 0 ]; then
#  pkill -e ino.out
#fi

pkill -e ino.out || exit 0