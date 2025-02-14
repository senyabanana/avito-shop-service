FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o avito-shop-service ./cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/avito-shop-service .
COPY .env .env

EXPOSE 8080

CMD ["./avito-shop-service"]