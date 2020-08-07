#!/bin/sh

if [[ $UID !=0 ]]; then
    echo "Please run this script with sudo:"
    echo "sudo $0 $*"
    exit 1
fi

dl=`goget -dir ~/Downloads -show true`

rm -rf /usr/local/go
tar -C /usr/local $dl
