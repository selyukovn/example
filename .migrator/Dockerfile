# ######################################################################################################################
# Сборка librdkafka
# ######################################################################################################################

# `pip install confluent-kafka` требует наличия `librdkafka v2.13.2+`,
# но `apk add librdkafka-dev` ставит 2.12 даже с `aok update`.
# Поэтому нужная версия собирается в отдельном шаге.

FROM python:3-alpine AS librdkafka

WORKDIR /librdkafka

RUN wget https://github.com/confluentinc/librdkafka/archive/refs/tags/v2.14.0.tar.gz && \
    tar -xzf v2.14.0.tar.gz && rm v2.14.0.tar.gz && \
    cd ./librdkafka-2.14.0 && \
    apk add g++ bash make && \
    ./configure --prefix=/librdkafka/compiled && make && make install && \
    cd ../ && rm -r ./librdkafka-2.14.0

# ######################################################################################################################
# Мигратор
# ######################################################################################################################

FROM python:3-alpine

# <confluent-kafka> -- `pip install confluent-kafka` не сработает без C-библиотек.
RUN apk add gcc libc-dev
COPY --from=librdkafka /librdkafka/compiled /usr
# </confluent-kafka>

WORKDIR /migrations

COPY up.py ./
COPY down.py ./
COPY entrypoint.sh ./
RUN chmod +x ./entrypoint.sh

VOLUME /migrations/scripts
VOLUME /migrations/state
VOLUME /migrations/packages

ENTRYPOINT ["/bin/sh", "/migrations/entrypoint.sh"]
CMD ["help"]

# ######################################################################################################################
