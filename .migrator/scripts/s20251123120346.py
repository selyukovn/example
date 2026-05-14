import mysql.connector
import os

class s20251123120346:
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
                    ALTER TABLE session
                        ADD COLUMN is_closed TINYINT(1) NOT NULL DEFAULT FALSE AFTER expire_at;
                """)
                cursor.execute("""
                    SET SESSION sql_safe_updates = 0;
                    UPDATE session
                    SET is_closed = 1
                    WHERE closed_at IS NOT NULL
                """)
                cursor.execute("""
                    CREATE INDEX i_going_expire ON session(is_closed, expire_at);
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
                cursor.execute("ALTER TABLE session DROP INDEX i_going_expire")
                cursor.execute("ALTER TABLE session DROP COLUMN is_closed")
