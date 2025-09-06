FROM nginx:latest

RUN rm /etc/nginx/conf.d/default.conf

COPY nginx.conf /etc/nginx/nginx.conf

COPY static/ /usr/share/nginx/html/static/

COPY html/ /usr/share/nginx/html/

EXPOSE 82
