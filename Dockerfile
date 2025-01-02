FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

ADD vizz /usr/bin/vizz

EXPOSE 8080

//ADD .env /.env

CMD ["vizz"]
