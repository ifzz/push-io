#!/usr/bin/env bash

sudo docker build -t controller:latest .

sudo docker images

sudo docker run -v /var/log:/log --name controller --restart=always -itd controller:latest

sudo docker ps
