#!/usr/bin/env bash

sudo docker build -t broker:latest .

sudo docker images

sudo docker run -v /var/log:/log --name broker --restart=always -itd -p 80:8080 broker:latest

sudo docker ps
