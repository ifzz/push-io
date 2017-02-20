#!/usr/bin/env bash

if [ -f ./mian ]
then
    rm ./main
fi

rm ./main

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
