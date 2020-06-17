#!/bin/bash

# connect 15 times from zookeeper clients

for I in 1 2 3 .. 15
do
    ~/zookeeper_ensemble/zookeeper-3.4.6_2/bin/zkCli.sh -server localhost:2182 <<EOF
get /t1
quit
EOF
sleep 2
done
