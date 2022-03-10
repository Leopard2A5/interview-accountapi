FROM golang:1.18rc1-alpine

RUN apk --update add build-base curl

WORKDIR /app

COPY entrypoint.sh ./
COPY *.go *.mod *.sum ./

ENTRYPOINT ./entrypoint.sh
