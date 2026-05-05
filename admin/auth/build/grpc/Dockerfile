
# Этапы сборки:
# - inner_build     -- (внутренний этап) создание бинарника.
# - inner_debugger  -- (внутренний этап) установка дебаггера.
# - debug           -- запуск с дебаггером Delve (например, для подключения к GoLand).
# - default         -- запуск в обычном режиме (по умолчанию, т.к. описан последним).

# Итоговые образы основаны на "чистых" alpine, т.к. golang-образы занимают много места и нужны только для компиляции.
# golang-latest = 900+mb
# golang-alpine = 200+mb
# alpine = 20+mb

# ----------------------------------------------------------------------------------------------------------------------
# Создание бинарника
# ----------------------------------------------------------------------------------------------------------------------

FROM golang:1.22.0-alpine AS inner_build

# Для компиляции с `confluent-kafka-go` нужен `gcc`.
# И поскольку образ основан на alpine, то вместо `libc-dev` следует использовать `musl-dev`.
RUN apk add gcc musl-dev

WORKDIR /app

# <Кеширование зависимостей в нижних слоях>
COPY app/go.mod app/go.sum ./
RUN go mod download
# </>

# <Компиляция>
COPY app ./
# https://github.com/confluentinc/confluent-kafka-go/tree/v2.12.0
# "... CGO_ENABLED must NOT be set to 0 since the Go client is based on the C library librdkafka. ..."
# "... If you are building for Alpine Linux (musl), -tags musl must be specified. ..."
# "... If the bundled librdkafka build is not supported on your platform,
# or you need a librdkafka with GSSAPI/Kerberos support, you must install librdkafka manually ..."
RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -o ./exec ./cmd/grpc.main.go
# </>

# ----------------------------------------------------------------------------------------------------------------------
# Установка дебаггера Delve
# ----------------------------------------------------------------------------------------------------------------------

FROM golang:1.22.0-alpine AS inner_debugger

RUN go install github.com/go-delve/delve/cmd/dlv@v1.25.2

# ----------------------------------------------------------------------------------------------------------------------
# Запуск с дебаггером
# ----------------------------------------------------------------------------------------------------------------------

FROM alpine:3.22.1 AS debug

WORKDIR /app

COPY --from=inner_debugger /go/bin/dlv   /app/dlv
COPY --from=inner_build /app/exec        /app/exec

VOLUME /state

# логика
EXPOSE 50051
# мониторинг
EXPOSE 8080
# дебаггер
EXPOSE 9999

ENV MYSQL_USER=""
ENV MYSQL_PASSWORD=""
ENV MYSQL_DB=""
ENV MYSQL_HOST_MASTER=""
ENV MYSQL_HOST_REPLICA=""
ENV API_GRPC_APIKEY=""
ENV SERVICE_CFM_API_GRPC_BASEURL=""
ENV SERVICE_CFM_API_GRPC_APIKEY=""

CMD ["./dlv", "--listen=:9999", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/exec"]

# ----------------------------------------------------------------------------------------------------------------------
# Запуск в обычном режиме (по умолчанию)
# ----------------------------------------------------------------------------------------------------------------------

FROM alpine:3.22.1 AS default

WORKDIR /app

COPY --from=inner_build /app/exec /app/exec

VOLUME /state

# логика
EXPOSE 50051
# мониторинг
EXPOSE 8080

ENV MYSQL_USER=""
ENV MYSQL_PASSWORD=""
ENV MYSQL_DB=""
ENV MYSQL_HOST_MASTER=""
ENV MYSQL_HOST_REPLICA=""
ENV API_GRPC_APIKEY=""
ENV SERVICE_CFM_API_GRPC_BASEURL=""
ENV SERVICE_CFM_API_GRPC_APIKEY=""

CMD ["./exec"]
