#!/bin/sh
set -e

if [ "${INIT_REPLICATOR_PASSWORD}" = "" ]; then
  echo "INIT_REPLICATOR_PASSWORD не может быть пуст!"
  exit 1
fi
if [ "${INIT_REPLICATION_HOST}" = "" ]; then
  echo "INIT_REPLICATION_HOST не может быть пуст!"
  exit 1
fi
if [ "${INIT_REPLICATION_SLOT}" = "" ]; then
  echo "INIT_REPLICATION_SLOT не может быть пуст!"
  exit 1
fi

# --

if [ -d "${PGDATA}" ] && [ ! -f "${PGDATA}"/PG_VERSION ]; then
  echo "Ожидание мастера (${INIT_REPLICATION_HOST}) ..."
  until pg_isready -h "${INIT_REPLICATION_HOST}" -U "replicator"; do
    echo "..."
    sleep 2
  done

  echo "Инициализация реплики..."
  export PGPASSWORD="${INIT_REPLICATOR_PASSWORD}"
  pg_basebackup \
    -d "host=${INIT_REPLICATION_HOST} user=replicator dbname=replication" \
    -D "${PGDATA}" -Fp \
    --checkpoint=fast -Xstream \
    -R -S "${INIT_REPLICATION_SLOT}" \
    --progress
  unset PGPASSWORD
  echo "Инициализация реплики завершена!"
fi

# --

exec /usr/local/bin/docker-entrypoint.sh "$@"
