# https://habr.com/ru/companies/doubletapp/articles/736602/
# "... Есть официальный экспортер от самой Nginx Inc.
# Для стандартной версии Nginx можно получить информацию об активных соединениях и количестве обработанных запросов,
# но нам нужно больше информации, как минимум HTTP-коды ответов и время обработки запросов.
# Официальный экспортер может собирать эти данные из API Nginx, но она доступна только в версии Nginx Plus.
# Для бесплатного Nginx также существует альтернативный экспортер — martin-helmich/prometheus-nginxlog-exporter.
# ..."
FROM ghcr.io/martin-helmich/prometheus-nginxlog-exporter/exporter:v1

COPY build/metrics-exporter/config.yml /etc/config.yml

VOLUME /input

CMD ["--config-file", "/etc/config.yml"]
