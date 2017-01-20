#!/usr/bin/env bash

for host in '54.223.146.80'
do
    echo $host

    scp ./main  ubuntu@$host:/home/ubuntu/push-io/controller/
done

for host in '10.35.68.215'
do
    echo $host

    scp ./main  gf@$host:/home/gf/push-io/controller/
done
