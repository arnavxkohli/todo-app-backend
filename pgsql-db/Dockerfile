FROM postgres:16

ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD 1234
ENV POSTGRES_DB todo-db

RUN rm -rf /var/lib/postgresql/data/*

COPY init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432
