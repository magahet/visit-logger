# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /httplogger

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /httplogger /httplogger

EXPOSE 5000

USER nonroot:nonroot

ENTRYPOINT ["/httplogger"]