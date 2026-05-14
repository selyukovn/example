import mysql.connector
import os

class s20251001083632:
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
                    CREATE TABLE account (
                        id BINARY(36) NOT NULL,
                        email VARBINARY(255) NOT NULL,
                        is_active BOOL NOT NULL,
                        deactivated_at DATETIME NULL DEFAULT NULL,
                        ip_whitelist_json VARBINARY(255) NOT NULL DEFAULT '[]',
                        created_at DATETIME NOT NULL,
                        updated_at DATETIME NOT NULL,
                        PRIMARY KEY(id),
                        UNIQUE KEY u_email (email)
                    )
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
                cursor.execute("DROP TABLE account")
