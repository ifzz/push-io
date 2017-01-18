#!/usr/bin/env bash

for host in '54.223.22.37'
do
    echo $host

    scp ./main  ubuntu@$host:/home/ubuntu/push-io/broker/
done
