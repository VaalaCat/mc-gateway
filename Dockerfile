FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache tzdata git
COPY go.mod go.sum ./

RUN CGO_ENABLED=0 go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o mcgw main.go

FROM alpine:latest
WORKDIR /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache tzdata git && mkdir -p /data
COPY --from=builder /app/mcgw .

EXPOSE 9000
EXPOSE 9001

ENV DB_PATH /data/data.db
ENTRYPOINT ["/app/mcgw"]