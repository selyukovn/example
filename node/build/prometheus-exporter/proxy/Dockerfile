FROM golang:1.22.0-alpine AS inner_build

WORKDIR /app

COPY build/prometheus-exporter/proxy/main.go ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o ./exec ./main.go

FROM alpine:3.22.1 AS default

ENV PORT_FOR_PROMETHEUS=9100
ENV PORT_FOR_PUSH_FROM_LOCAL=9101

WORKDIR /app

COPY --from=inner_build /app/exec /app/exec

CMD ["./exec"]
