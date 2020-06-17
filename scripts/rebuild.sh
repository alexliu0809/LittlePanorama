#!/bin/bash

# use this script to clean the current zookeeper build
# build again
# then create a new three-server ensemble

./stop_all.sh

rm -rf ~/zookeeper_ensemble/*

ant clean
ant package
cp zoo.cfg build/zookeeper-3.4.6/conf/
cd build/
cd zookeeper-3.4.6/
mkdir data
cd ..

echo "1" > zookeeper-3.4.6/data/myid
cp -r zookeeper-3.4.6/ ~/zookeeper_ensemble/
mv ~/zookeeper_ensemble/zookeeper-3.4.6 ~/zookeeper_ensemble/zookeeper-3.4.6_1
cp ../zoo1.cfg ~/zookeeper_ensemble/zookeeper-3.4.6_1/conf/zoo.cfg


echo "2" > zookeeper-3.4.6/data/myid
cp -r zookeeper-3.4.6/ ~/zookeeper_ensemble/
mv ~/zookeeper_ensemble/zookeeper-3.4.6 ~/zookeeper_ensemble/zookeeper-3.4.6_2
cp ../zoo2.cfg ~/zookeeper_ensemble/zookeeper-3.4.6_2/conf/zoo.cfg


echo "3" > zookeeper-3.4.6/data/myid
cp -r zookeeper-3.4.6/ ~/zookeeper_ensemble/
mv ~/zookeeper_ensemble/zookeeper-3.4.6 ~/zookeeper_ensemble/zookeeper-3.4.6_3
cp ../zoo3.cfg ~/zookeeper_ensemble/zookeeper-3.4.6_3/conf/zoo.cfg

# auto cleaning
cd ..
ant clean



#cp zoo.cfg build/zookeeper-3.4.6/conf/
#cd build/zookeeper-3.4.6/
#mkdir data

# bin/zkServer.sh start
# bin/zkCli.sh
# https://zookeeper.apache.org/doc/r3.3.3/zookeeperStarted.html
# bin/zkServer.sh stop
# sudo netstat -plnt
