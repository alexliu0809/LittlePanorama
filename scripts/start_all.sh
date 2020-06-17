#!/bin/bash

# start all services

rm deephealth.db

nohup myhview-server -addr localhost:6688 pano0 > hview-server.out 2>&1 &

cd ~/zookeeper_ensemble/zookeeper-3.4.6_1
./bin/zkServer.sh start

cd ~/zookeeper_ensemble/zookeeper-3.4.6_2
./bin/zkServer.sh start

cd ~/zookeeper_ensemble/zookeeper-3.4.6_3
./bin/zkServer.sh start


