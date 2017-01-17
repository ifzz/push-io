#!/usr/bin/env bash

for host in '54.223.146.80'
do
    echo $host

    scp ./main  ubuntu@$host:/home/ubuntu/push-io/controller/
done
