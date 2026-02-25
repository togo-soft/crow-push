FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && CGO_ENABLED=0 go build -o /plugin

FROM alpine:latest
WORKDIR /app
COPY --from=builder /plugin /usr/local/bin/plugin
ENTRYPOINT ["/usr/local/bin/plugin"]
