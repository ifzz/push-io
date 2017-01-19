#!/usr/bin/env bash

for host in '54.223.22.37' '119.29.20.165'
do
    echo $host

    scp ./main  ubuntu@$host:/home/ubuntu/push-io/monitor/
done
