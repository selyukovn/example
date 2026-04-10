#!/bin/sh
set -e

vType=${KAFKA_PROCESS_ROLES}

if [ "$vType" = "broker" ]; then
  kafka-metadata-quorum --bootstrap-server=localhost:9092 describe --status > /dev/null
elif [ "$vType" = "controller" ]; then
  kafka-metadata-quorum --bootstrap-controller=localhost:9093 describe --status > /dev/null
else
  # broker + controller
  kafka-metadata-quorum --bootstrap-server=localhost:9092 describe --status > /dev/null
  kafka-metadata-quorum --bootstrap-controller=localhost:9093 describe --status > /dev/null
fi
