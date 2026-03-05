#!/bin/sh

healthcheck_str="HEALTHCHECK:OK:$(date)"

mysql -u docker_healthcheck_user -e "SELECT '${healthcheck_str}' AS healthcheck;" \
  | grep "${healthcheck_str}" \
  > /tmp/healthcheck

if [ "$(cat /tmp/healthcheck)" != "${healthcheck_str}" ]; then
  # Вывод можно будет найти в docker inspect
  echo "ОШИБКА: \"$(cat /tmp/healthcheck)\" != \"${healthcheck_str}\""
  exit 1
fi
