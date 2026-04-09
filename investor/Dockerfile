FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY investor/go.mod investor/go.sum ./

RUN go mod download

COPY investor/ ./investor/
COPY plantform/ ./plantform/

RUN CGO_ENABLED=0 GOOS=linux go build -o /investor ./investor/cmd/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

COPY --from=builder /investor /investor

ENV POLL_INTERVAL=2s
ENV SYMBOLS=SBER,GAZP,MOEX
ENV LOG_LEVEL=info

ENTRYPOINT ["/investor"]
