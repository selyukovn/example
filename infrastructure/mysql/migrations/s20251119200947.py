import mysql.connector
import os

class s20251119200947:
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
                    CREATE TABLE session (
                        id BINARY(36) NOT NULL,
                        account_id BINARY(36) NOT NULL,
                        sign_in_request_id BINARY(36) NOT NULL,
                        initial_client_user_agent VARBINARY(255) NOT NULL,
                        initial_client_ip VARCHAR(45) NOT NULL,
                        initiated_at DATETIME NOT NULL,
                        expire_at DATETIME NOT NULL,
                        closed_at DATETIME NULL DEFAULT NULL,
                        created_at DATETIME NOT NULL,
                        updated_at DATETIME NOT NULL,
                        PRIMARY KEY(id),
                        UNIQUE KEY u_sign_in_request (sign_in_request_id)
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
                cursor.execute("DROP TABLE session")
