# syntax=docker/dockerfile:1

FROM golang:1.21.6-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth-service ./cmd/main.go

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /bin/auth-service /app/

ARG VERSION

ENV SERVICE_NAME="auth-service" \
    SERVICE_VERSION=${VERSION}

EXPOSE 8081

CMD ["/app/auth-service"]
