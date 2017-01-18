#!/usr/bin/env bash

sudo docker build -t gateway:latest .

sudo docker images

sudo docker run --restart=always -itd --name gateway -p 80:3000 gateway:latest

sudo docker ps
