#!/bin/bash

trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT
set -e

go run ./server &
sleep 1
go run ./client


