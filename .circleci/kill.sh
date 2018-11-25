#!/bin/sh

pkill -e ino.out > /dev/null 2>&1

if [ $? = 0 ]; then
  pkill -e ino.out
fi