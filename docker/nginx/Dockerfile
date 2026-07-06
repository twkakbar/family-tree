FROM nginx:1.27-alpine

RUN rm /etc/nginx/conf.d/default.conf \
    && mkdir -p /var/www/uploads

COPY nginx.conf /etc/nginx/nginx.conf
COPY conf.d/ /etc/nginx/conf.d/
