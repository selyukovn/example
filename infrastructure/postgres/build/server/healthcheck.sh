#!/bin/sh

if [ "$(pg_isready -U postgres | grep "accepting connections")" = "" ]; then
  # Вывод можно будет найти в docker inspect
  echo "ОШИБКА: healthcheck провалился"
  exit 1
fi
