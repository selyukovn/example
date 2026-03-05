FROM golang:1.22.0-alpine AS inner_build

WORKDIR /app

COPY build/prometheus-exporter/local-pusher/main.go ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o ./exec ./main.go

FROM alpine:3.22.1 AS default

ENV PORT_EXPORTER=9100
ENV PORT_PROXY_PUSH=9101
ENV INTERVAL_SECONDS=5

WORKDIR /app

COPY --from=inner_build /app/exec /app/exec

CMD ["./exec"]
