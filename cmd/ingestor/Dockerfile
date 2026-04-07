FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ingestor ./cmd/ingestor

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/ingestor .

ENV KAFKA_BROKER=localhost:9092
ENV FINNHUB_API_KEY=

EXPOSE 8080

CMD ["./ingestor"]
