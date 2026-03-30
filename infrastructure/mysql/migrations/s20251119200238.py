import mysql.connector
import os

class s20251119200238:
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
                    CREATE TABLE action_request (
                        id BINARY(36) NOT NULL,
                        type TINYINT UNSIGNED NOT NULL,
                        account_id BINARY(36) NOT NULL,
                        confirmation_id BINARY(36) NOT NULL,
                        requested_at DATETIME NOT NULL,
                        created_at DATETIME NOT NULL,
                        updated_at DATETIME NOT NULL,
                        PRIMARY KEY(id),
                        UNIQUE KEY u_confirmation (confirmation_id)
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
                cursor.execute("DROP TABLE action_request")
