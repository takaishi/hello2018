#!/bin/bash
echo ">> START"
echo "FOO = ${FOO}"
#sleep 30

if [ "${FOO}" = "false" ]; then
  echo "Failed..."
  exit 1
fi
echo "Hello World!!!1"
echo ">> FINISH"
