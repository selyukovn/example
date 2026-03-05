import mysql.connector
import os

class s20251121235532:
    def up(self):
        with mysql.connector.connect(
            host=os.environ['MYSQL_HOST'],
            port=3306,
            user=os.environ['MYSQL_USER'],
            password=os.environ['MYSQL_PASSWORD'],
            database=os.environ['MYSQL_DB_AUTH'],
        ) as connection:
            with connection.cursor() as cursor:
                cursor.execute("""
                    CREATE TABLE event (
                        id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
                        occurred_at DATETIME NOT NULL,
                        type VARBINARY(255) NOT NULL,
                        version TINYINT UNSIGNED NOT NULL,
                        extra_account_email VARBINARY(255) NULL DEFAULT NULL,
                        extra_account_id BINARY(36) NULL DEFAULT NULL,
                        extra_account_ip_whitelist_json VARBINARY(255) NULL DEFAULT NULL,
                        extra_session_id BINARY(36) NULL DEFAULT NULL,
                        created_at DATETIME NOT NULL,
                        PRIMARY KEY (id)
                    )
                """)

    def down(self):
        with mysql.connector.connect(
            host=os.environ['MYSQL_HOST'],
            port=3306,
            user=os.environ['MYSQL_USER'],
            password=os.environ['MYSQL_PASSWORD'],
            database=os.environ['MYSQL_DB_AUTH'],
        ) as connection:
            with connection.cursor() as cursor:
                cursor.execute("DROP TABLE event")
