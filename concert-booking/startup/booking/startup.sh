#!/bin/bash

set -m

trap exithandler CHLD

exithandler() {
  if [ -n "${READY}" ]; then
    if kill -0 ${LINKERDPID} && kill -0 ${BOOKINGPID}; then
      echo "SIGCHLD received, but all child processes alive, doing nothing"
    else
      echo "SIGCHLD received, killing child processes"
      kill ${LINKERDPID} || true
      kill ${BOOKINGPID} || true
      sleep 3
      kill -9 ${LINKERDPID} || true
      kill -9 ${BOOKINGPID} || true
      exit 1
    fi
  fi
}

pushd /app/
/app/linkerd-1.3.6-exec /linkerd.yml &
LINKERDPID=$!
popd

pushd /app/
/app/booking &
BOOKINGPID=$!
popd

READY=true

wait
