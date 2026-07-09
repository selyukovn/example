#!/bin/sh
set -e

if [ "${INIT_POSTGRES_PASSWORD}" = "" ]; then
  echo "INIT_POSTGRES_PASSWORD не может быть пустым!"
  exit 1
fi
if [ "${INIT_PROMETHEUS_EXPORTER_PASSWORD}" = "" ]; then
  echo "INIT_PROMETHEUS_EXPORTER_PASSWORD не может быть пустым!"
  exit 1
fi
if [ "${INIT_MIGRATOR_PASSWORD}" = "" ]; then
  echo "INIT_MIGRATOR_PASSWORD не может быть пустым!"
  exit 1
fi
if [ "${INIT_REPLICATOR_PASSWORD}" = "" ]; then
  echo "INIT_REPLICATOR_PASSWORD не может быть пустым!"
  exit 1
fi
if [ "${INIT_REPLICATION_SLOT_1}" = "" ]; then
  echo "INIT_REPLICATION_SLOT_1 не может быть пустым!"
  exit 1
fi

# --

# Имя "POSTGRES_PASSWORD" не говорит об использовании только на этапе инициализации,
# поэтому введен более мнемонический псевдоним.
export POSTGRES_PASSWORD="${INIT_POSTGRES_PASSWORD}"

# --

exec /usr/local/bin/docker-entrypoint.sh "$@"
