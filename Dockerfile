FROM golang:1.12-alpine AS build

LABEL maintainer "enm10k <enm10k@gmail.com>"

RUN apk add --update ca-certificates git openssh
ADD . /go/src/github.com/shiguredo/ayame
RUN cd /go/src/github.com/shiguredo/ayame && \
    GO111MODULE=on CGO_ENABLED=0 go build -o /usr/bin/ayame

EXPOSE 3000
WORKDIR /go/src/github.com/shiguredo/ayame
ENTRYPOINT ["/usr/bin/ayame", "-addr=0.0.0.0:3000"]
