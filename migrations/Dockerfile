FROM mariadb:latest

EXPOSE 3306/tcp

COPY . /docker-entrypoint-initdb.d

ENV MYSQL_RANDOM_ROOT_PASSWORD true
ENV MYSQL_USER bast
ENV MYSQL_PASSWORD bast
ENV MYSQL_DATABASE bast

