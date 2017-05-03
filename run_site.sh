#!/bin/bash

PS_LINE=`sudo ps -A | grep fire`
if [ "$PS_LINE" ]; then
  echo "Process found"
  PID=`echo "$PS_LINE" | cut -d' ' -f2`
  sudo kill $PID && echo "Process killed"
else
  echo "Process not found"
fi
go build && echo "Site built. Starting..." && (sudo ./fire-phillips-data &)
