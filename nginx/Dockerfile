FROM nginx:latest

COPY ./web/* /usr/share/nginx/html/
COPY ./jsmodel /usr/share/nginx/html/jsmodel
COPY ./nginx.conf /etc/nginx/nginx.conf
VOLUME /upload
