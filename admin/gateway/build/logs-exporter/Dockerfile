FROM grafana/alloy:v1.12.2

# http://loki:3100
ARG LOKI_BASEURL

COPY build/logs-exporter/config.alloy                   /etc/alloy/config.alloy
RUN sed -i -e "s|__LOKI_BASEURL__|${LOKI_BASEURL}|g"    /etc/alloy/config.alloy

VOLUME /input/http
VOLUME /input/kafcon
VOLUME /state

EXPOSE 12345

CMD ["run", "/etc/alloy/config.alloy", "--storage.path=/state", "--server.http.listen-addr=0.0.0.0:12345"]
