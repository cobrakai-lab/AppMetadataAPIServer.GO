# syntax=docker/dockerfile:1
FROM golang:latest

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN go build -o main .
EXPOSE 8888
CMD ["/app/main"]