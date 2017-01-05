#!/usr/bin/env bash

sudo docker run --restart=always -itd --name collector collector:latest

sudo docker ps
