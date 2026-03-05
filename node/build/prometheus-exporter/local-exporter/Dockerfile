FROM prom/node-exporter:v1.10.2

ENV PORT=9100
ENV EXTRA_FLAGS=""

# Exec-форма не интерпретирует переменные,
# поэтому обязательна "shell-форма + exec" для использования того же PID=1 для корректной обработки сигналов.
# 127.0.0.1 -- обязательно, чтобы метрики не были доступны на других интерфейсах.
# EXTRA_FLAGS нужен тут в таком виде, поскольку в "shell-форме + exec"
# передача флагов через "command" в docker-compose.yml файле не работает.
ENTRYPOINT exec /bin/node_exporter --path.rootfs=/host --web.listen-address="127.0.0.1:${PORT}" ${EXTRA_FLAGS}
