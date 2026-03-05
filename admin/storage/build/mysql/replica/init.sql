
CHANGE REPLICATION SOURCE TO
    SOURCE_HOST='__REPLICA_SOURCE_HOST__',
    SOURCE_PORT=__REPLICA_SOURCE_PORT__,

    # Файл и позиция бинлога не указаны, т.к. это скрипт инициализации новой реплики,
    # т.е. наполняться она будет с начальной позиции по умолчанию.

    SOURCE_USER='docker_replication_user',
    SOURCE_PASSWORD='docker_replication_password',

    SOURCE_SSL=0,
    GET_SOURCE_PUBLIC_KEY=1;
