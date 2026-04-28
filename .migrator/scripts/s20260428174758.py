import mysql.connector
import os

class s20260428174758:
    def up(self):
        with mysql.connector.connect(
            host=os.environ['MYSQL_HOST'],
            port=3306,
            user=os.environ['MYSQL_USER'],
            password=os.environ['MYSQL_PASSWORD'],
            database="auth",
        ) as connection:
            with connection.cursor() as cursor:
                cursor.execute("""
                    ALTER TABLE event
                        ADD COLUMN outbox_group_id VARBINARY(255) NULL DEFAULT NULL,
                        ADD COLUMN outbox_operation_id VARBINARY(255) NULL DEFAULT NULL
                """)
                cursor.execute("""
                    UPDATE event
                    SET
                        outbox_group_id = extra_account_id,
                        outbox_operation_id = (SELECT UUID_TO_BIN(UUID()))
                    WHERE id > 0
                """)
                cursor.execute("""
                    ALTER TABLE event
                        MODIFY outbox_group_id VARBINARY(255) NOT NULL
                        COMMENT "идентификатор группы событий / ключ партицирования для механизма распространения событий",
                        MODIFY outbox_operation_id VARBINARY(255) NOT NULL
                        COMMENT "идентификатор операции, инициировавшей событие"
                """)

    def down(self):
        with mysql.connector.connect(
            host=os.environ['MYSQL_HOST'],
            port=3306,
            user=os.environ['MYSQL_USER'],
            password=os.environ['MYSQL_PASSWORD'],
            database="auth",
        ) as connection:
            with connection.cursor() as cursor:
                cursor.execute("""
                    ALTER TABLE event
                        DROP COLUMN outbox_group_id,
                        DROP COLUMN outbox_operation_id
                """)
