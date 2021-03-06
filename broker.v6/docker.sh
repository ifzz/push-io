#!/usr/bin/env bash

sudo docker build -t docker.gf.com.cn/broker:latest .

sudo docker images

sudo docker run -v /var/log:/log --name broker --restart=always -itd -p 8080:8080 docker.gf.com.cn/broker:latest

sudo docker ps
