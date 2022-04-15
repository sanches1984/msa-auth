FROM golang:1.16-alpine as builder
WORKDIR /app
ENV GO111MODULE=on

ARG version
ENV AUTH_APP_VERSION=$version

COPY . .
RUN go mod download
RUN go build -o auth cmd/main.go

FROM alpine
RUN apk update && \
    adduser -D -H -h /app auth && \
    mkdir -p /app/migrations  && \
    chown -R auth:auth /app
WORKDIR /app
USER auth

COPY --chown=auth --from=builder /app/auth /app
COPY --chown=auth --from=builder /app/migrations /app/migrations
COPY --chown=auth --from=builder /app/.env /app

EXPOSE 5000
CMD ["/app/auth"]
