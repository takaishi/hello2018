#!/bin/bash
set -x

echo ">> START"
apk update
apk upgrade
apk add libcurl
apk add curl
apk add jq
apk add base64 
FOO=$(curl -XGET ${CONSUL_HOST}:8500/v1/kv/foo| jq -r ".[]|.Value" | base64 -d)
echo "FOO = ${FOO}"
#sleep 30

if [ "${FOO}" = "false" ]; then
  echo "Failed..."
  exit 1
fi
echo "Hello World!!!1"
echo ">> FINISH"
