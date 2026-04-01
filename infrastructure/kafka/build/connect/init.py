import urllib.request
import urllib.parse
import json
import os

allBrokersString = os.environ["CONNECT_BOOTSTRAP_SERVERS"]
brokersCount = max(1, len(allBrokersString.split(",")))

# https://habr.com/ru/companies/deliveryclub/articles/529484/
# https://debezium.io/documentation/reference/3.5/connectors/mysql.html
# https://docs.confluent.io/platform/current/connect/references/restapi.html#connectors

def mysqlEventProducer(
    connectorName,
    dbTables,
    topicRedirectFromRegexp,
    topicRedirectTo,
    partitionsCount,
    partitionKeyingFields,
):
    req = urllib.request.Request(
        "http://localhost:8083/connectors",
        method="POST",
        headers={"Content-Type": "application/json"},
        data=json.dumps({
            "name": connectorName,
            "config": {
                "connector.class": "io.debezium.connector.mysql.MySqlConnector",

                # ---- mysql ----
                "database.server.id": 100501,
                "database.hostname": os.environ["MYSQL_HOST"],
                "database.port": os.environ["MYSQL_PORT"],
                "database.user": os.environ["MYSQL_USER"],
                "database.password": os.environ["MYSQL_PASSWORD"],
                "table.include.list": dbTables,
                # ---- /mysql ----

                # ---- schema.history ----
                # databases -- однозначно true, т.к. в теории разные базы -- это разные сервера, а значит и коннекторы.
                # tables -- коннектор играет роль продюсера событий, поэтому использование других таблиц не ожидается.
                "schema.history.internal.store.only.captured.databases.ddl": True,
                "schema.history.internal.store.only.captured.tables.ddl": True,
                # ---- /schema.history ----

                # ---- topics ----
                "schema.history.internal.kafka.bootstrap.servers": allBrokersString,

                # По умолчанию topic.creation.enable = true, поэтому будут созданы топики:
                #
                # - `connectorName` с DDL:
                #   SET character_set_server=utf8mb4, collation_server=utf8mb4_0900_ai_ci
                #   DROP TABLE IF EXISTS `auth`.`event`
                #   DROP DATABASE IF EXISTS `auth`
                #   CREATE DATABASE `auth` CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci
                #   USE `auth`
                #   CREATE TABLE `event` ...
                #
                # - технический топик для внутренних дел -- `connectorName + "-history"`:
                #   SET character_set_server=utf8mb4, collation_server=utf8mb4_0900_ai_ci
                #   DROP TABLE IF EXISTS `auth`.`event`
                #   CREATE TABLE `event` ...
                #
                # - по каждой табличке из dbTables -- `connectorName + "." + {dbName.tableName}`:
                #   ... данные ...
                "schema.history.internal.kafka.topic": connectorName + "-history",
                "topic.prefix": connectorName,

                # https://debezium.io/documentation/reference/3.5/configuration/topic-auto-create-config.html
                # По умолчанию топики будут созданы только с 1 партицией без реплик.
                # В топике `connectorName` будет DDL -- он не нужен для event-producer'ования.
                # Но для топика с данными это снижает надежность и лишает возможности параллельно обрабатывать события,
                # т.е. в рамках одного сервиса-приемника может быть запущен только один консюмер на группу
                # (по факту больше может быть и не нужно, но перспектива дальнейшего развития должна быть открыта).
                #
                # Внимание!
                # Использование нескольких партиций обязует таблицы предоставлять ключ партицирования,
                # по которому будет вычисляться конкретная партиция, в которую будет отправлено конкретное событие.
                # Это необходимо для того, чтобы события из определенных групп (например, по пользователю)
                # принимались в том же порядке, в котором они были записаны в таблицу событий,
                # а также чтобы события из разных групп могли быть обработаны параллельно.
                # Технически можно собрать такой ключ прямо здесь (см. документацию), но это место легко забыть,
                # а последствия непоследовательной обработки событий тихо распространятся на всю систему.
                # Поэтому лучше иметь специальную колонку в таблице событий, которая будет видна при внесении изменений.
                # Таблица событий ведется именно для их распространения, поэтому никакие зависимости не будут нарушены.
                #
                # Суть `event-producer`-коннектора в использовании одной таблички
                # (или нескольких табличек-частей логически целой таблицы, например, events_2025, events_2026, ...),
                # поэтому достаточно одной группы параметров создания топиков,
                # тем более в данном случае из-за "редиректа" все события перенаправляются в единственный топик.
                # Однако, debezium требует обязательно указать `replication.factor` и `partitions` для группы `default`.
                # Топик `connectorName` с DDL будет создан по default-группе, поскольку не войдет в вышеописанную.
                "topic.creation.default.replication.factor": brokersCount,
                "topic.creation.default.partitions": 1,
                "topic.creation.groups": topicRedirectTo,
                "topic.creation."+topicRedirectTo+".replication.factor": brokersCount,
                "topic.creation."+topicRedirectTo+".partitions": partitionsCount,
                "topic.creation."+topicRedirectTo+".include": topicRedirectTo,

                # Для наглядности, сокрытия реализации продюсера и эксперимента ради используется "публичный" топик.
                # Он будет создан вместо `connectorName + ".dbName.tableName"` топика(-ов?).
                # https://kafka.apache.org/42/kafka-connect/user-guide/#org.apache.kafka.connect.transforms.RegexRouter
                "transforms.RedirectToPublicTopic.type": "org.apache.kafka.connect.transforms.RegexRouter",
                "transforms.RedirectToPublicTopic.regex": topicRedirectFromRegexp,
                "transforms.RedirectToPublicTopic.replacement": topicRedirectTo,

                # https://debezium.io/documentation/reference/3.5/transformations/partition-routing.html
                "transforms.Partitioning.type": "io.debezium.transforms.partitions.PartitionRouting",
                "transforms.Partitioning.partition.payload.fields": partitionKeyingFields,
                "transforms.Partitioning.partition.topic.num": partitionsCount,
                "transforms.Partitioning.predicate": topicRedirectTo,
                "predicates": topicRedirectTo,
                "predicates."+topicRedirectTo+".type": "org.apache.kafka.connect.transforms.predicates.TopicNameMatches",
                "predicates."+topicRedirectTo+".pattern": topicRedirectTo,

                "transforms": "RedirectToPublicTopic,Partitioning",
                # ---- /topics ----
            }
        }).encode("utf-8"),
    )

    try:
        with urllib.request.urlopen(req) as response:
            result = json.loads(response.read().decode("utf-8"))
            print("Статус:", response.getcode())
            print("Ответ:", result)
    except urllib.error.HTTPError as e:
        error_details = e.read().decode("utf-8")
        print(f"Status Code: {e.code}")
        print(f"Reason: {e.reason}")
        print(f"Body: {error_details}")
        raise e

# ######################################################################################################################
# ВЫПОЛНЕНИЕ
# ######################################################################################################################


# ######################################################################################################################
