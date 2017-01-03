#!/usr/bin/env bash

for host in '54.223.124.84'
do
    echo $host

    scp ./main  ubuntu@$host:/home/ubuntu/push-io/monitor/
done
