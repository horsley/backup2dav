FROM golang:latest AS builder
WORKDIR /app
COPY *.go go.* /app/
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go mod tidy
RUN go build -o main

FROM alpine:latest

WORKDIR /
COPY --from=builder /app/main /
ENV CRON_SCHEDULE="3 0 * * *"
COPY setup_cron.sh /

CMD ["sh", "/setup_cron.sh"]