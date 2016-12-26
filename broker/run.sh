#!/usr/bin/env bash

sudo docker build -t console:latest .

sudo docker images

sudo docker run -v /var/log:/log --name console --restart=always -d -p 8080:8080 console:latest

sudo docker ps
