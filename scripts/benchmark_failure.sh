#!/bin/bash

# Scripts to benchmark failures (gary failures + crash failures)

# Rebuild everything

cd ~/zookeeper-release-3.4.6

################ First, Crash Leader ################
echo "Creating Testing Environment For Leader Failure"
## reset environment
./rebuild.sh >/dev/null 2>&1
sleep 15
./start_all.sh >/dev/null 2>&1
sleep 15

## Get pid to kill
for i in $(sudo netstat -tulpn | grep "127.0.0.1:3882" | tr " " "\n")
do
    if [[ "$i" == *"java"* ]]; then
        PID=$(echo $i | tr -dc '0-9')
        #echo $PID
        break
    fi
done

## set start time
start=$(date +%s%N)/1000000

## Kill PID
kill -9 $PID

## Keep pulling
while true
do
    r=$(hview-client -server 127.0.0.1:6688 get inference 2)
    #echo $r
    if [[ "$r" == *"UNHEALTHY"* ]]; then
        # save end, break
    end=$(date +%s%N)/1000000
    break
    fi 
done
echo "Time To Detect Leader Failure: $(($end-$start)) ms"
echo ""


################ Second, Crash Follower ################
echo "Creating Testing Environment For Follower Failure"
## reset environment
./rebuild.sh >/dev/null 2>&1
sleep 15
./start_all.sh >/dev/null 2>&1
sleep 15

## Get pid to kill
for i in $(sudo netstat -tulpn | grep "127.0.0.1:3881" | tr " " "\n")
do
    if [[ "$i" == *"java"* ]]; then
        PID=$(echo $i | tr -dc '0-9')
        #echo $PID
        break
    fi
done

## set start time
start=$(date +%s%N)/1000000

## Kill PID
kill -9 $PID

## Keep pulling
while true
do
    r=$(hview-client -server 127.0.0.1:6688 get inference 1)
    #echo $r
    if [[ "$r" == *"UNHEALTHY"* ]]; then
        # save end, break
    end=$(date +%s%N)/1000000
    break
    fi 
done
echo "Time To Detect Follower Failure: $(($end-$start)) ms"
echo ""



################ Third, Gray 1 ################
echo "Creating Testing Environment For Gray Failure 1"

## reset environment
./rebuild.sh >/dev/null 2>&1
sleep 15
./start_all.sh >/dev/null 2>&1
sleep 15

## set up necessary things
echo "Setting up nodes"
~/zookeeper_ensemble/zookeeper-3.4.6_2/bin/zkCli.sh -server localhost:2182 >/dev/null 2>&1 <<EOF
create /t1 t1
create /t2 t2
quit
EOF
sleep 1

## background script keep pulling
echo "Setting up regular connections from clients"
nohup ./keep_connect.sh >/dev/null 2>&1 &
echo $! > save_pid.txt

## Trigger gray1
echo "Triggering gray 1"
start=$(date +%s%N)/1000000
~/zookeeper_ensemble/zookeeper-3.4.6_2/bin/zkCli.sh -server localhost:2182 >/dev/null 2>&1 <<EOF
create /gray1 gray1
quit
EOF


## Keep pulling
while true
do
    r=$(hview-client -server 127.0.0.1:6688 get inference 2)
    #echo $r
    if [[ "$r" == *"UNHEALTHY"* ]]; then
        # save end, break
    end=$(date +%s%N)/1000000
    break
    fi 
done
echo "Time To Detect Gray Failure 1: $(($end-$start)) ms"
echo ""

# cleanup no hup
#kill -9 `cat save_pid.txt`
#rm save_pid.txt
sleep 5


################ Fourth, Gray 2 ################
echo "Creating Testing Environment For Gray Failure 2"
## reset environment
./rebuild.sh >/dev/null 2>&1
sleep 15
./start_all.sh >/dev/null 2>&1
sleep 15

## Trigger gray2
echo "Triggering gray 2"
start=$(date +%s%N)/1000000
~/zookeeper_ensemble/zookeeper-3.4.6_2/bin/zkCli.sh -server localhost:2182 >/dev/null 2>&1 <<EOF
create /gray2 gray2
quit
EOF

## Keep pulling
while true
do
    r=$(hview-client -server 127.0.0.1:6688 get inference 2)
    #echo $r
    if [[ "$r" == *"UNHEALTHY"* ]]; then
        # save end, break
    end=$(date +%s%N)/1000000
    break
    fi 
done
echo "Time To Detect Gray Failure 2: $(($end-$start)) ms"
echo ""


################ Fifth, Coming Back ################
echo "Creating Testing Environment For Coming Back"
## reset environment
./rebuild.sh >/dev/null 2>&1
sleep 15
./start_all.sh >/dev/null 2>&1
sleep 15

## set up necessary things
echo "Setting up nodes"
~/zookeeper_ensemble/zookeeper-3.4.6_2/bin/zkCli.sh -server localhost:2182 >/dev/null 2>&1 <<EOF
create /t1 t1
create /t2 t2
quit
EOF
sleep 1

## background script keep pulling
echo "Setting up regular connections from clients"
nohup ./keep_connect.sh >/dev/null 2>&1 &

## Trigger gray1
echo "Triggering gray 1"
start=$(date +%s%N)/1000000
~/zookeeper_ensemble/zookeeper-3.4.6_2/bin/zkCli.sh -server localhost:2182 >/dev/null 2>&1 <<EOF
create /gray1 gray1
quit
EOF


## Keep pulling
while true
do
    r=$(hview-client -server 127.0.0.1:6688 get inference 2)
    #echo $r
    if [[ "$r" == *"UNHEALTHY"* ]]; then
        # save end, break
    end=$(date +%s%N)/1000000
    break
    fi 
done
echo "Time To Detect Failure: $(($end-$start)) ms"


## Start Pulling Again
sleep 45
echo "Start Pulling Again From Clients"
nohup ./keep_connect_more.sh >/dev/null 2>&1 &

while true
do
    r=$(hview-client -server 127.0.0.1:6688 get inference 2)
    #echo $r
    if [[ "$r" == *"network: HEALTHY"* ]]; then
        # save end, break
    end2=$(date +%s%N)/1000000
    break
    fi 
done
echo "Time To Detect Revival: $(($end2-$start)) ms"
echo ""
./stop_all.sh >/dev/null 2>&1 
rm nohup.out >/dev/null 2>&1 
