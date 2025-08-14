FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /main .
COPY .env .env

EXPOSE 8080 2112

ENTRYPOINT ["/app/main"]
