#!/usr/bin/env bash

trap "rm tinykv;kill 0" EXIT

go build -o tinykv

./tinykv -c p1.yml &
./tinykv -c p2.yml &
./tinykv -c p3.yml &

wait
