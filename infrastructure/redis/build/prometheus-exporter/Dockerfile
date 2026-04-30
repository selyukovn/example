FROM oliver006/redis_exporter:v1.82.0-alpine

ENV REDIS_EXPORTER_REDIS_ONLY_METRICS = true

# https://hub.docker.com/r/oliver006/redis_exporter#prometheus-configuration-to-scrape-multiple-redis-hosts
CMD ["--redis.addr", ""]
