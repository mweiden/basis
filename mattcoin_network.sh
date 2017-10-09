#!/usr/bin/env bash
set -ex

PIDS=""
PORTS=`seq 5000 5002`

finish() {
  for pid in ${PIDS} ; do
    kill -9 $pid
  done
}

trap "finish" EXIT

for PORT in ${PORTS} ; do
  ./mattcoin -port=${PORT} &
  PIDS="$! $PIDS"
done

sleep 5

for CONTACT_PORT in ${PORTS} ; do
  for REGISTER_PORT in ${PORTS} ; do
    curl -XPOST \
      -d "{\"address\":\"localhost:${REGISTER_PORT}\"}" \
      -H 'Content-Type: application/json' \
      localhost:${CONTACT_PORT}/nodes/register
  done
done

while True ; do
  sleep 5
  for CONTACT_PORT in ${PORTS} ; do
    for REGISTER_PORT in ${PORTS} ; do
      curl -XGET localhost:${CONTACT_PORT}/nodes/resolve
    done
  done
done
