
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

WORKDIR /app

# <Кеширование зависимостей в нижних слоях>
COPY app/go.mod app/go.sum ./
RUN go mod download
# </>

# <Компиляция>
COPY app ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./exec ./cmd/http.main.go
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

# <Копирование файлов>
COPY --from=inner_build /app/exec       /app/exec
COPY app/cmd/http/bundles/auth/static   /app/static/auth
# дебагер
COPY --from=inner_debugger /go/bin/dlv /app/dlv
# </>

VOLUME /state

# логика
EXPOSE 80
# мониторинг
EXPOSE 8080
# дебаггер
EXPOSE 9999

ENV APP_NAME=""
ENV BASE_URL=""
ENV SESSION_COOKIE_NAME=""
ENV SERVICE_AUTH_API_GRPC_BASEURL=""
ENV SERVICE_AUTH_API_GRPC_APIKEY=""

CMD ["./dlv", "--listen=:9999", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/exec"]

# ----------------------------------------------------------------------------------------------------------------------
# Запуск в обычном режиме (по умолчанию)
# ----------------------------------------------------------------------------------------------------------------------

FROM alpine:3.22.1 AS default

WORKDIR /app

# <Копирование файлов>
COPY --from=inner_build /app/exec       /app/exec
COPY app/cmd/http/bundles/auth/static   /app/static/auth
# </>

VOLUME /state

# логика
EXPOSE 80
# мониторинг
EXPOSE 8080

ENV APP_NAME=""
ENV BASE_URL=""
ENV SESSION_COOKIE_NAME=""
ENV SERVICE_AUTH_API_GRPC_BASEURL=""
ENV SERVICE_AUTH_API_GRPC_APIKEY=""

CMD ["./exec"]
