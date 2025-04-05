FROM golang:1.23.0-alpine3.18 AS Build

WORKDIR /app
COPY . .
RUN go mod download

RUN go build -o /auth-service main.go

FROM alpine:latest

WORKDIR /
COPY --from=Build /auth-service /auth-service

ENTRYPOINT [ "/auth-service" ]