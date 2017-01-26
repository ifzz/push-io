#!/usr/bin/env bash

START_TIME=$SECONDS

for i in {1..1000}
do
    curl -XPOST 'http://push-it.gf.com.cn/api/v1/notification'  -H 'Content-Type:application/json' --data '{"appId":"gftrader","appKey":"1163CFFD87155CD634CBD3DA9F53D","topic": "mike","message":{"payload":{"广发证券":"涨停", "招商证券":"跌停"}}}'

    printf "\n\n"
done

ELAPSED_TIME=$(($SECONDS - $START_TIME))

printf "\n\n"
echo "$(($ELAPSED_TIME/60)) min $(($ELAPSED_TIME%60)) sec"
