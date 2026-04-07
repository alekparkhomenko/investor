FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o alert-engine ./cmd/alert-engine

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/alert-engine .

ENV KAFKA_BROKER=localhost:9092
ENV REDIS_ADDR=localhost:6379

EXPOSE 8081

CMD ["./alert-engine"]
