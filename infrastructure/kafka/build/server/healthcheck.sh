#!/bin/sh
set -e

vType=${KAFKA_PROCESS_ROLES}

if [ "$vType" = "broker" ]; then
  kafka-metadata-quorum --bootstrap-server=localhost:9092 describe --status > /dev/null
elif [ "$vType" = "controller" ]; then
  kafka-metadata-quorum --bootstrap-controller=localhost:9092 describe --status > /dev/null
else
  echo "ОШИБКА! HEALTHCHECK не определен для узла комбинированного типа!"
  exit 1
fi
