#! /usr/bin/env bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <name>"
    exit 1
fi

migrate create -ext sql -dir ./migrations -seq $1
