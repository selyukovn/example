FROM golang:1.22.0-alpine

ARG SERVICE_NAME

# <grpc>
WORKDIR /protoc
RUN apk add curl && \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v30.2/protoc-30.2-linux-x86_64.zip && \
    unzip protoc-30.2-linux-x86_64.zip && rm protoc-30.2-linux-x86_64.zip && \
    # плагины protoc'а для генерации кода на Go
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.7 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1 && \
    # полезные proto-файлы
    curl -LO https://github.com/googleapis/googleapis/archive/refs/tags/common-protos-1_3_1.zip && \
    unzip common-protos-1_3_1.zip && rm common-protos-1_3_1.zip && \
    mv googleapis-common-protos-1_3_1/google/* include/google/ && rm -r googleapis-common-protos-1_3_1
ENV PATH="${PATH}:/protoc/bin"
# Скрипт для удобства запуска генерации
RUN echo "protoc \\" > /protoc/generate_grpc_api.sh && \
    echo "--go_out=/app/cmd/grpc \\" >> /protoc/generate_grpc_api.sh && \
    echo "--go-grpc_out=/app/cmd/grpc \\" >> /protoc/generate_grpc_api.sh && \
    echo "--proto_path=/app/cmd/grpc/pb \\" >> /protoc/generate_grpc_api.sh && \
    echo "/app/cmd/grpc/pb/${SERVICE_NAME}.proto" >> /protoc/generate_grpc_api.sh && \
    chmod +x /protoc/generate_grpc_api.sh
# </grpc>

# <open-api>
WORKDIR /openapi
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1
RUN echo "rm -f /app/cmd/http/bundles/\$1/openapi/*.openapi.go" >> /openapi/gen.sh \
    # разбивка по файлам работает не как ожидается: типы response'ы будут в файле с сервером, а не в файле с типами,
    # поэтому весь код записывается в один файл.
    && echo "oapi-codegen -generate=types,std-http,strict-server -o /app/cmd/http/bundles/\$1/openapi/server.openapi.go /app/cmd/http/bundles/\$1/openapi/spec/openapi.yml" >> /openapi/gen.sh \
    && echo "oapi-codegen -generate=types,client -o /app/cmd/http/bundles/\$1/openapi/client.openapi.go /app/cmd/http/bundles/\$1/openapi/spec/openapi.yml" >> /openapi/gen.sh \
    && chmod +x /openapi/gen.sh
# </open-api>

# <app>
WORKDIR /app
VOLUME /app

# Для установки некоторых пакетов `go get` требует установленного `git`.
RUN apk add git

# Некоторые команды (например, `go mod tidy` или `go get`)
# требуют права записи как для внутренних директорий (например, /.cache/go-build), так и для подключенных в VOLUME /app.
# Для использования внутренних файлов нужно запускаться под `root`, ведь у "обычного" пользователя не будет нужных прав.
# Но при таком запуске придется каждый раз менять владельца файлов, созданных в VOLUME /app, что очень неудобно.
ARG USER
RUN chown -R ${USER} /go \
    && mkdir -p /.cache/go-build && chown -R ${USER} /.cache/go-build
# </app>

# --

# Игнорируем предупреждение "JSONArgsRecommended",
# т.к. CMD требуется прописать не по форме, а корректная обработка сигналов ОС в данном случае неважна.
# https://docs.docker.com/reference/build-checks/json-args-recommended/#explicitly-specify-the-shell
SHELL ["/bin/sh", "-c"]

CMD echo "Готово! Жду -- подключайтесь..." \
    # Чтобы контейнер был всегда запущен -- для предотвращения перезапуска при watch-ребилде любого другого контейнера.
    && touch /tmp/tail && tail -F /tmp/tail > /dev/null
