# syntax = docker/dockerfile:1.2
ARG BASE=golang:1.18.9
FROM ${BASE} AS builder

WORKDIR /top

COPY . .

ENV GOPROXY=https://goproxy.cn,direct
#RUN go build -o xrelayer main.go

CMD ["/bin/bash", "-c", "./build.sh"]