#!/bin/sh

echo "Building Windows command line binary"
GOOS=windows GOARCH=amd64 go build -o getgo.exe

echo "Building Linux command line binary"
GOOS=linux GOARCH=amd64 go build -o getgo

echo "Building Linux Raspberry Pi command line binary"
GOOS=linux GOARCH=arm go build -o getgo-arm

echo "Building Linux GUI interface"
GOOS=linux GOARCH=amd64 go build -tags gui,fyne -o getgo-gui

echo "Building Windows GUI interface"
CC=x86_64-w64-mingw32-gcc GOOS=windows CGO_ENABLED=1 go build -ldflags -H=windowsgui -tags fyne,gui -o getgo-win.exe 
