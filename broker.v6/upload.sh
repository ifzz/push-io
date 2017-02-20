#!/usr/bin/env bash

for host in '54.223.22.37'
do
    echo $host

    scp ./main  ubuntu@$host:/home/ubuntu/push-io/broker/
done


#for host in '10.2.68.215' '10.35.68.215'
#do
#    echo $host
#
#    scp ./main  gf@$host:/home/gf/push-io/broker/
#done
