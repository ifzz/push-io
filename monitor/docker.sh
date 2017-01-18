#!/usr/bin/env bash

sudo docker build -t monitor:latest .

sudo docker images

sudo docker run -v /var/log:/log --name monitor --restart=always -itd monitor:latest

sudo docker ps
