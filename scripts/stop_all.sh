#!/bin/bash

# stop_all services


cd ~/zookeeper_ensemble/zookeeper-3.4.6_1
./bin/zkServer.sh stop

cd ~/zookeeper_ensemble/zookeeper-3.4.6_2
./bin/zkServer.sh stop

cd ~/zookeeper_ensemble/zookeeper-3.4.6_3
./bin/zkServer.sh stop


pkill java
pkill hview-server
pkill myhview-server

rm deephealth.db*
#kill -9 `cat save_pid.txt`
#rm save_pid.txt
