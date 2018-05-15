#!/bin/bash
echo ">> START"
apk update
apk add curl
FOO=$(curl -XGET ${CONSUL_HOST}:8500/v1/kv/foo)
echo "FOO = ${FOO}"
#sleep 30

if [ "${FOO}" = "false" ]; then
  echo "Failed..."
  exit 1
fi
echo "Hello World!!!1"
echo ">> FINISH"
