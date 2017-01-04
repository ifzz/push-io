#!/usr/bin/env bash

sudo docker run --restart=always -itd --name gateway -p 3000:3000 gateway:latest

sudo docker ps
