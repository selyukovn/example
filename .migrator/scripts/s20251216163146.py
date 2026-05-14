import mysql.connector
import os

class s20251216163146:
    def up(self):
        with mysql.connector.connect(
            host=os.environ['MYSQL_HOST'],
            port=3306,
            user=os.environ['MYSQL_USER'],
            password=os.environ['MYSQL_PASSWORD'],
            database="confirmation",
        ) as connection:
            with connection.cursor() as cursor:
                cursor.execute("""
                    CREATE TABLE event (
                        id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
                        occurred_at DATETIME NOT NULL,
                        type VARBINARY(255) NOT NULL,
                        version TINYINT UNSIGNED NOT NULL,
                        extra_confirmation_id BINARY(36) NULL DEFAULT NULL,
                        extra_finish_type TINYINT UNSIGNED NULL DEFAULT NULL,
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
            database="confirmation",
        ) as connection:
            with connection.cursor() as cursor:
                cursor.execute("""DROP TABLE event""")
