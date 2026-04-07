FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o notifier ./cmd/notifier

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/notifier .

ENV KAFKA_BROKER=localhost:9092
ENV TELEGRAM_TOKEN=

EXPOSE 8082

CMD ["./notifier"]
