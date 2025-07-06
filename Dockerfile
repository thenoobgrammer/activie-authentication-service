FROM golang:1.24.4-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/authentication-service ./cmd/main.go

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /bin/authentication-service /app/

EXPOSE 8081

CMD ["/app/authentication-service"]
