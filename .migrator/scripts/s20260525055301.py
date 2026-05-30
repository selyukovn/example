import mysql.connector
import os

class s20260525055301:
    def up(self):
        with mysql.connector.connect(
            host=os.environ['MYSQL_HOST'],
            port=3306,
            user=os.environ['MYSQL_USER'],
            password=os.environ['MYSQL_PASSWORD'],
            database="gateway",
        ) as connection:
            with connection.cursor() as cursor:
                cursor.execute("""
                    CREATE TABLE sys_dlq (
                        id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
                        topic VARBINARY(255) NOT NULL,
                        group_id VARBINARY(255) NOT NULL,
                        m_key BLOB NOT NULL,
                        m_value BLOB NOT NULL,
                        m_partition INT NOT NULL,
                        m_offset BIGINT NOT NULL,
                        m_metadata VARBINARY(255) NULL DEFAULT NULL,
                        m_headers_keys_json BLOB NOT NULL,
                        m_headers_values_json BLOB NOT NULL,
                        m_timestamp DATETIME NOT NULL,
                        m_timestamp_type TINYINT NOT NULL,
                        created_at DATETIME NOT NULL,
                        PRIMARY KEY (id),
                        INDEX i_topic_group (topic, group_id)
                    )
                """)

    def down(self):
        with mysql.connector.connect(
            host=os.environ['MYSQL_HOST'],
            port=3306,
            user=os.environ['MYSQL_USER'],
            password=os.environ['MYSQL_PASSWORD'],
            database="gateway",
        ) as connection:
            with connection.cursor() as cursor:
                cursor.execute("DROP TABLE sys_dlq")
