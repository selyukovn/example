FROM python:3-alpine

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
