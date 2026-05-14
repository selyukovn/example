FROM alpine:3.22.1

# "Как приручить драCRONа"
# https://habr.com/ru/companies/redmadrobot/articles/305364/

# Сам logrotate
# https://man7.org/linux/man-pages/man8/logrotate.8.html

WORKDIR /logrotate

# --------------------
# Установка
# --------------------

RUN apk add --no-cache logrotate

# Удаляем дефолтные конфиги и кронтаб, т.к. нужно будет перенаправить вывод в докер
RUN rm /etc/logrotate.conf \
    && rm -r /etc/logrotate.d \
    && rm /etc/periodic/*/logrotate

# Канал, а не файл, т.к. нужен только "сквозной" вывод
RUN mkfifo -m 0666 output.fifo

# Выносим файл состояния в вольюм.
RUN mkdir -m 0500 state.d && touch state.d/state # права  тут неважны, т.к. будет переопределены при подключении
VOLUME /logrotate/state.d

# Скрипт для cron'а
RUN echo "#!/bin/sh"                                                                          >> logrotate_cron \
 && echo "/usr/sbin/logrotate" \
    "--state=\"/logrotate/state.d/state\"" \
    "--log=\"/logrotate/output.fifo\"" \
    "/logrotate/logrotate.conf"                                                               >> logrotate_cron \
 && echo "EXITVALUE=\$?"                                                                      >> logrotate_cron \
 && echo "if [ \$EXITVALUE != 0 ]; then"                                                      >> logrotate_cron \
 && echo "    /usr/bin/logger -t logrotate \"ALERT exited abnormally with [\$EXITVALUE]\""    >> logrotate_cron \
 && echo "fi"                                                                                 >> logrotate_cron \
 && echo "exit 0"                                                                             >> logrotate_cron \
 && chmod +x logrotate_cron \
 && mv logrotate_cron /etc/periodic/daily/

# Скрипт запуска контейнера
RUN echo "#!/usr/bin/env sh"                            >> entrypoint \
 && echo "crond"                                        >> entrypoint \
 && echo "trap \"killall crond; exit\" SIGINT SIGTERM"  >> entrypoint \
 && echo "tail -f /logrotate/output.fifo & wait \$!"    >> entrypoint \
 && chmod 0555 entrypoint

# --------------------
# Наведение
# --------------------

COPY build/logrotate/logrotate.conf ./
RUN chmod 0444 logrotate.conf

VOLUME /app/admin/auth/grpc
VOLUME /app/admin/auth/sch
VOLUME /app/admin/cfm/grpc
VOLUME /app/admin/cfm/sch
VOLUME /app/admin/gateway/http
VOLUME /app/admin/gateway/kafcon
VOLUME /app/infrastructure/mysql
VOLUME /app/infrastructure/redis

# --------------------
# Запуск
# --------------------

ENTRYPOINT ["/logrotate/entrypoint"]
