FROM danielqsj/kafka-exporter:v1.9.0

ENV BROKER_A=""
ENV BROKER_B=""
ENV BROKER_C=""

# Exec-форма не интерпретирует переменные,
# поэтому обязательна "shell-форма + exec" для использования того же PID=1 для корректной обработки сигналов.
ENTRYPOINT exec /bin/kafka_exporter --kafka.server=${BROKER_A} --kafka.server=${BROKER_B} --kafka.server=${BROKER_C}
