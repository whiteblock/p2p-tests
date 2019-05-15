#!/bin/bash

export NO_PRETTY=1

#(size,name)
log(){
	for ((i=0;i<$1;i++))
	do
		wb get log $i > ".data/${2}/node${i}.log"
	done
}

wait_for_results(){
	sleep 2m
}


SENDERS=50
NODES=100
SIZE=200
CONNS=10
INTERVAL=7812
psr="floodsub"

#SERIES 1 Control
#a
wb build -b libp2p-test -n $NODES -m 0 -c 0 -y -o"senders=$SENDERS" -o"payloadSize=$SIZE" -o"connections=$CONNS" -o"interval=$INTERVAL" -o="pubsubRouter=$psr"

wait_for_results

log $NODES series1a

#b
wb build -b libp2p-test -n $NODES -m 0 -c 0 -y -o"senders=$SENDERS" -o"payloadSize=$SIZE" -o"connections=$CONNS" -o"interval=$INTERVAL" -o="pubsubRouter=$psr"

wait_for_results

log $NODES series1b

#c
wb build -b libp2p-test -n $NODES -m 0 -c 0 -y -o"senders=$SENDERS" -o"payloadSize=$SIZE" -o"connections=$CONNS" -o"interval=$INTERVAL" -o="pubsubRouter=$psr"

wait_for_results

log $NODES series1c


#Series 2: Message Size

SIZE=500
#a
wb build -b libp2p-test -n $NODES -m 0 -c 0 -y -o"senders=$SENDERS" -o"payloadSize=$SIZE" -o"connections=$CONNS" -o"interval=$INTERVAL" -o="pubsubRouter=$psr"

wait_for_results

log $NODES series2a

SIZE=500000
#b
wb build -b libp2p-test -n $NODES -m 0 -c 0 -y -o"senders=$SENDERS" -o"payloadSize=$SIZE" -o"connections=$CONNS" -o"interval=$INTERVAL" -o="pubsubRouter=$psr"

wait_for_results

log $NODES series2b

SIZE=500000000
#c
wb build -b libp2p-test -n $NODES -m 0 -c 0 -y -o"senders=$SENDERS" -o"payloadSize=$SIZE" -o"connections=$CONNS" -o"interval=$INTERVAL" -o="pubsubRouter=$psr"

wait_for_results

log $NODES series2c

#Series 3: Senders

SIZE=200
NODES=100
SENDERS=10
#a
wb build -b libp2p-test -n $NODES -m 0 -c 0 -y -o"senders=$SENDERS" -o"payloadSize=$SIZE" -o"connections=$CONNS" -o"interval=$INTERVAL" -o="pubsubRouter=$psr"

wait_for_results

log $NODES series3a

SENDERS=40

#b
wb build -b libp2p-test -n $NODES -m 0 -c 0 -y -o"senders=$SENDERS" -o"payloadSize=$SIZE" -o"connections=$CONNS" -o"interval=$INTERVAL" -o="pubsubRouter=$psr"

wait_for_results

log $NODES series3b

SENDERS=90

#c
wb build -b libp2p-test -n $NODES -m 0 -c 0 -y -o"senders=$SENDERS" -o"payloadSize=$SIZE" -o"connections=$CONNS" -o"interval=$INTERVAL" -o="pubsubRouter=$psr"

wait_for_results

log $NODES series3c
