#!/usr/bin/env bash

for host in '54.222.243.29'
do
    echo $host

    scp ./main  ubuntu@$host:/home/ubuntu/push-io/broker/
done
