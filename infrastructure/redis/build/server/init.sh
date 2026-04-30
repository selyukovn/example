#!/bin/sh
set -e

# ----------------------------------------------------------------------------------------------------------------------
# ACL
# ----------------------------------------------------------------------------------------------------------------------

# https://redis.io/docs/latest/operate/oss_and_stack/management/security/acl/

# Вместо default для явного подключения и необходимости использования default для healthcheck (см. ниже).
echo 'ACL SETUSER root on >initial_root_password sanitize-payload ~* &* +@all' | redis-cli

# Для cache-сервера миграции, конечно, смысла не имеют, но не хочется разделять из-за этого скрипт инициализации.
# Для db-сервера миграции все же могут иметь смысл, хоть в точности и сложно предположить что-то конкретное.
# Пусть начальный набор правил совпадает с app_user -- по мере использования будет видно, что нужно добавить.
echo 'ACL SETUSER migration_user on >initial_migration_password sanitize-payload ~* &* ' \
  '-@all +@connection +@read +@write -@dangerous' \
  | redis-cli

echo 'ACL SETUSER app_user on >initial_app_password sanitize-payload ~* &* ' \
  '-@all +@connection +@read +@write -@dangerous' \
  | redis-cli

# https://github.com/oliver006/redis_exporter/tree/v1.82.0#authenticating-with-redis
echo 'ACL SETUSER prometheus_exporter_user on >initial_prometheus_exporter_password' \
  '-@all +@connection +memory -readonly +strlen +config|get +xinfo +pfcount -quit +zcard +type +xlen ' \
  '-readwrite -command +client -wait +scard +llen +hlen +get +eval +slowlog ' \
  '+cluster|info +cluster|slots +cluster|nodes -hello -echo +info +latency +scan -reset -auth -asking ' \
  | redis-cli

# Для выполнения инициализации требуется, чтобы сервер был запущен,
# а значит healthcheck-скрипт должен работать и до инициализации (см. launcher.sh),
# т.е. до создания специализированного healthcheck-пользователя.
# Такое возможно только с использованием пользователя default --
# поэтому вместо создания healthcheck-пользователя передаем эту роль пользователю default, отбирая все остальные права.
# Делать это нужно в конце, поскольку инициализацию тоже выполняет default, а ACL отключает права сразу.
# Пароль для healthcheck-проверок (ping) по сути не нужен,
# но для того, чтобы разрешить подключаться к redis из других контейнеров (например, для экспорта метрик),
# необходимо либо отключить protected mode, либо задать пароль default-пользователю.
# "... Redis is running in protected mode because protected mode is enabled and no password is set for the default user.
# In this mode connections are only accepted from the loopback interface. ..."
# Выбран вариант с использованием пароля для default пользователя, поскольку так и так уже вносятся изменения в ACL.
# К тому же значение пароля "healthcheck" подчеркивает роль пользователя,
# а само наличие пароля предотвращает неявное подключение от имении default-пользователя.
echo 'ACL SETUSER default on >healthcheck -@all +ping' | redis-cli

# В `redis.conf` должен быть указан `aclfile` иначе результат инициализации не сохранится для следующего запуска.
# Можно бы было воспользоваться "не багом, а фичей" и выполнять такую "инициализацию" при каждом запуске контейнера,
# но такой подход отличается от инициализации других компонентов инфраструктуры -- лучше таки поддерживать однообразие.
# Выполнять уже нужно от root'а, поскольку default-пользователь уже лишен права на ACL SAVE.
redis-cli --no-auth-warning -u redis://root:initial_root_password@127.0.0.1:6379 ACL SAVE

# ----------------------------------------------------------------------------------------------------------------------
