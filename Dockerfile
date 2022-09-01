FROM alpine:latest

RUN mkdir /app

COPY ./bin/app_linux ./app/api

CMD ["/app/api"]