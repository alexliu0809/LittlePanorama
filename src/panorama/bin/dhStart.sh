#!/bin/bash
DHBIN="${BASH_SOURCE-$0}"
DHBIN="$(dirname "${DHBIN}")"
DHBINDIR="$(cd "${DHBIN}"; pwd)"

if [ $# -eq 0 ]; then
  config="$DHBINDIR/../hs.cfg"
elif [ $# -eq 1 ]; then
  config=$1
else
  echo "Usage: $0 CONFIG"
  exit 1
fi
if [ ! -f $config ]; then
	echo "Could not find config file $config"
	exit 1
fi

if [ -f deephealth.pid ]; then
  echo "Deep health server process has already started. Stop it first."
  exit 0
fi

hview-server -config $config > deephealth.out 2>&1 &
dh_pid=$!
sleep 1
if ps -p$dh_pid > /dev/null; then
  echo $dh_pid > deephealth.pid
  echo "Deep health server started with PID $dh_pid"
else
  echo "Deep health server has exited"
fi
