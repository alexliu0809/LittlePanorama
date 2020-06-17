#!/bin/bash

# Connect 3 times from zookeeper client

sleep 3
for i in 1 2 3
do
    ~/zookeeper_ensemble/zookeeper-3.4.6_2/bin/zkCli.sh -server localhost:2182 <<EOF
get /t1
quit
EOF
done
