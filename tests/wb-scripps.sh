#!/bin/bash

export NO_PRETTY=1

n=20

wb build -b libp2p-test -n $n -m 0 -c 0 -y

sleep 2m

for ((i=0;i<$n;i++))
do
    wb get log $i > "/home/daniel/p2ptestdata/n${i}log.txt"
done