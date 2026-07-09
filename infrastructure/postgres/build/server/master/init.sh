#!/bin/sh
set -e

# SQL без подсветки синтаксиса не очень удобен в чтении,
# но напрямую sql-файлы использовать не удастся,
# поскольку инициализация требует подключения к разным бд,
# а множество отдельных шаблонов только усложнит процесс сборки.

sqlFile="/tmp/init.sql"

cat << 'SQL_FILE' > "${sqlFile}"
-- ---------------------------------------------------------------------------------------------------------------------
-- PROMETHEUS EXPORTER
-- ---------------------------------------------------------------------------------------------------------------------

CREATE USER prometheus_exporter WITH
    -- В документации к prometheus-экспортеру нет указания ограничивать кол-во соединений с сервером postgres
    -- для избежания лишней нагрузки от мониторинговых скрапов.
    -- https://grafana.com/oss/prometheus/exporters/postgres-exporter/?tab=installation
    --
    -- Однако, кажется, что смысл в таком ограничении все же есть -- потому сделано по аналогии с mysql.
    -- https://grafana.com/oss/prometheus/exporters/mysql-exporter/?tab=installation
    --
    CONNECTION LIMIT 3

    PASSWORD '__PROMETHEUS_EXPORTER_PASSWORD__';

-- Права предоставлены согласно доке экспортера:
-- https://github.com/prometheus-community/postgres_exporter?tab=readme-ov-file#running-as-non-superuser
--
GRANT pg_monitor TO prometheus_exporter;
GRANT pg_read_all_stats TO prometheus_exporter;

-- ---------------------------------------------------------------------------------------------------------------------
-- MIGRATION
-- ---------------------------------------------------------------------------------------------------------------------

CREATE USER migrator WITH
    -- Вероятно, нет смысла в количестве разрешенных соединений > 1.
    --
    -- https://www.postgresql.org/docs/17/sql-createrole.html
    -- "... The CONNECTION LIMIT option is only enforced approximately;
    -- if two new sessions start at about the same time when just one connection “slot” remains for the role,
    -- it is possible that both will fail. Also, the limit is never enforced for superusers. ..."
    --
    -- Но миграции выполняются последовательно, поэтому проблем с использованием CONNECTION LIMIT 1 не ожидается.
    --
    CONNECTION LIMIT 1

    PASSWORD '__MIGRATOR_PASSWORD__';

-- ---------------------------------------------------------------------------------------------------------------------
-- REPLICATION
-- ---------------------------------------------------------------------------------------------------------------------

-- https://www.postgresql.org/docs/17/warm-standby.html#STREAMING-REPLICATION-AUTHENTICATION
-- "... It is recommended to create a dedicated user account with REPLICATION and LOGIN privileges for replication. ..."
--
CREATE USER replicator WITH
    REPLICATION
    PASSWORD '__REPLICATOR_PASSWORD__';

SELECT * FROM pg_create_physical_replication_slot('__REPLICATION_SLOT_1__');

-- ---------------------------------------------------------------------------------------------------------------------
SQL_FILE

sed -i -e "s|__PROMETHEUS_EXPORTER_PASSWORD__|${INIT_PROMETHEUS_EXPORTER_PASSWORD}|g"   "${sqlFile}"
sed -i -e "s|__MIGRATOR_PASSWORD__|${INIT_MIGRATOR_PASSWORD}|g"                         "${sqlFile}"
sed -i -e "s|__REPLICATOR_PASSWORD__|${INIT_REPLICATOR_PASSWORD}|g"                     "${sqlFile}"
sed -i -e "s|__REPLICATION_SLOT_1__|${INIT_REPLICATION_SLOT_1}|g"                       "${sqlFile}"

psql -U postgres < "${sqlFile}"
