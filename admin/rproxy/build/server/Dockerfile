FROM nginx:1.29-alpine

# http://bff:8080
ARG BFF_PROXY_PASS

# <Конфиг>
COPY build/server/nginx.conf                                /etc/nginx/nginx.conf
RUN sed -i -e "s|__BFF_PROXY_PASS__|${BFF_PROXY_PASS}|g"    /etc/nginx/nginx.conf
# </>

# <Копируем в образ статику -- см. nginx.conf директиву root и location'ы>
COPY app/   /var/www/static
# </>

# <Eсли конфиг кривой, error_log в nginx.conf не сработает -- запись пойдет в /var/log/nginx по умолчанию.
# Поэтому симлинк обязателен в любом случае.
# Предварительно папку надо удалить, иначе будет "/var/log/nginx/state -> /state" вместо "/var/log/nginx -> /state".>
RUN rm -r /var/log/nginx && ln -sf /state/logs /var/log/nginx
VOLUME /state/logs
# При запуске с user-директивой в docker-compose, требует доступ к "/var/cache/nginx".
RUN rm -r /var/cache/nginx && ln -sf /state/cache /var/cache/nginx
VOLUME /state/cache
VOLUME /state/run
# </>

EXPOSE 80
