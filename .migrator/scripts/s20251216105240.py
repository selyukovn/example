import mysql.connector
import os

class s20251216105240:
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
                    CREATE TABLE confirmation (
                        id BINARY(36) NOT NULL,
                        email VARBINARY(255) NOT NULL,
                        expire_at DATETIME NOT NULL,
                        finished_at DATETIME NULL DEFAULT NULL,
                        finish_type TINYINT UNSIGNED NOT NULL,
                        fails_made TINYINT UNSIGNED NOT NULL,
                        created_at DATETIME NOT NULL,
                        updated_at DATETIME NOT NULL,
                        PRIMARY KEY (id)
                    )
                """)
                cursor.execute("""
                    CREATE TABLE confirmation_request (
                        confirmation_id BINARY(36) NOT NULL,
                        number TINYINT UNSIGNED NOT NULL,
                        code_hash VARBINARY(255) NOT NULL,
                        requested_at DATETIME NOT NULL,
                        PRIMARY KEY (confirmation_id, number)
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
                cursor.execute("""DROP TABLE confirmation_request""")
                cursor.execute("""DROP TABLE confirmation""")
