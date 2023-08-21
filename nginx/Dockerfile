FROM nginx:latest

COPY ./web/* /usr/share/nginx/html/
COPY ./jsmodel /usr/share/nginx/html/jsmodel
COPY ./nginx.conf /etc/nginx/nginx.conf

VOLUME /usr/share/nginx/html/images
