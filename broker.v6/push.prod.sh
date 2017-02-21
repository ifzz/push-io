#!/usr/bin/env bash

sudo docker build -t docker.gf.com.cn/broker:latest .

sudo docker images

sudo docker push docker.gf.com.cn/broker:latest

rm ./config.json

rm ./key.json
