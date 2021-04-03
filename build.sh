#!/bin/sh

echo "Building Windows command line binary"
GOOS=windows GOARCH=amd64 go build -o getgo.exe

echo "Building Linux command line binary"
GOOS=linux GOARCH=amd64 go build -o getgo

echo "Building Linux Raspberry Pi command line binary"
GOOS=linux GOARCH=arm go build -o getgo-arm
