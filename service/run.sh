#!/bin/bash
nohup /opt/dcontrol -p 888 > /opt/dcontrol.log 2>&1 &
nohup ./dcontrol -p 8085 > ./dcontrol.log 2>&1 &


# 找出包含关键字 "monitor" 的进程ID
pids=$(ps aux | grep -v grep | grep dcontrol | awk '{print $2}')
if [ -n "$pids" ]; then
    echo "找到包含关键字 'dcontrol' 的进程，将被杀死：$pids"
    kill -9 $pids
fi