#!/bin/bash

# create an ensemble of three servers

./stop_all.sh
rm -rf ~/zookeeper_ensemble/*

cd build

cp -r zookeeper-3.4.6/ ~/zookeeper_ensemble/
mv ~/zookeeper_ensemble/zookeeper-3.4.6 ~/zookeeper_ensemble/zookeeper-3.4.6_1
cp ../zoo1.cfg ~/zookeeper_ensemble/zookeeper-3.4.6_1/conf/zoo.cfg

cp -r zookeeper-3.4.6/ ~/zookeeper_ensemble/
mv ~/zookeeper_ensemble/zookeeper-3.4.6 ~/zookeeper_ensemble/zookeeper-3.4.6_2
cp ../zoo2.cfg ~/zookeeper_ensemble/zookeeper-3.4.6_2/conf/zoo.cfg

cp -r zookeeper-3.4.6/ ~/zookeeper_ensemble/
mv ~/zookeeper_ensemble/zookeeper-3.4.6 ~/zookeeper_ensemble/zookeeper-3.4.6_3
cp ../zoo3.cfg ~/zookeeper_ensemble/zookeeper-3.4.6_3/conf/zoo.cfg

cd ..
