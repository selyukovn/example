FROM golang:1.22.0-alpine

WORKDIR /ide

VOLUME ./admin/auth/go.mod
VOLUME ./admin/auth/go.sum

VOLUME ./admin/cfm/go.mod
VOLUME ./admin/cfm/go.sum

VOLUME ./admin/front/go.mod
VOLUME ./admin/front/go.sum

VOLUME ./admin/gateway/go.mod
VOLUME ./admin/gateway/go.sum

VOLUME ./path
VOLUME ./sdk

# Игнорируем предупреждение "JSONArgsRecommended",
# т.к. CMD требуется прописать не по форме, а корректная обработка сигналов ОС в данном случае неважна.
# https://docs.docker.com/reference/build-checks/json-args-recommended/#explicitly-specify-the-shell
SHELL ["/bin/sh", "-c"]

CMD echo $(date +"%Y/%m/%d %H:%M:%S") ": Установка пакетов для \"admin/auth\"..." \
    && cd /ide/admin/auth && go mod download \
    # --
    && echo $(date +"%Y/%m/%d %H:%M:%S") ": Установка пакетов для \"admin/cfm\"..." \
    && cd /ide/admin/cfm && go mod download \
    # --
    && echo $(date +"%Y/%m/%d %H:%M:%S") ": Установка пакетов для \"admin/front\"..." \
    && cd /ide/admin/front && go mod download \
    # --
    && echo $(date +"%Y/%m/%d %H:%M:%S") ": Установка пакетов для \"admin/gateway\"..." \
    && cd /ide/admin/gateway && go mod download \
    # ----
    && echo $(date +"%Y/%m/%d %H:%M:%S") ": Передача всех пакетов IDE..." \
    && chmod -Rf 0775 /ide/path \
    && rm -rf /ide/sdk/* /ide/path/* \
    && cp -r /usr/local/go/* /ide/sdk && cp -r /go/* /ide/path \
    && chmod -R 0775 /ide/path \
    # ----
    && echo $(date +"%Y/%m/%d %H:%M:%S") ": Готово!" \
    # Чтобы контейнер был всегда запущен -- для предотвращения перезапуска при watch-ребилде любого другого контейнера.
    && touch /tmp/tail && tail -F /tmp/tail > /dev/null
