FROM grafana/alloy:v1.12.2

# http://loki:3100
ARG LOKI_BASEURL

COPY build/logs-exporter/config.alloy                   /etc/alloy/config.alloy
RUN sed -i -e "s|__LOKI_BASEURL__|${LOKI_BASEURL}|g"    /etc/alloy/config.alloy

VOLUME /input
VOLUME /state

EXPOSE 12345

# 0.0.0.0 -- поскольку интерфейс 127.0.0.1 не доступен снаружи контейнера.
CMD ["run", "/etc/alloy/config.alloy", "--storage.path=/state", "--server.http.listen-addr=0.0.0.0:12345"]
