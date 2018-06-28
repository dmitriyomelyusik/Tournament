FROM golang

COPY ./bin /go/tournament/

ENV DBNAME=postgres PGPASS=password PGUSER=postgres PGHOST=127.0.0.1 SSLMODE=disable

CMD tournament/game

EXPOSE 8080
