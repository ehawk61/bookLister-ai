FROM postgres:17

ENV POSTGRES_DB=booklister
ENV POSTGRES_USER=booklister
ENV POSTGRES_PASSWORD=booklister

COPY db/bookLister.pg.sql /docker-entrypoint-initdb.d/

EXPOSE 5432
