FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.work and module files for both modules
COPY go.work go.work.sum ./
COPY investor/go.mod investor/go.sum ./
COPY plantform/go.mod plantform/go.sum ./

# Copy source code
COPY investor/ ./investor/
COPY plantform/ ./plantform/

# Copy swagger docs
COPY investor/internal/http/docs/swagger.json ./docs/swagger.json

# Build will download dependencies automatically
RUN CGO_ENABLED=0 GOOS=linux go build -o /investor ./investor/cmd/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

COPY --from=builder /investor /investor
COPY --from=builder /app/docs /docs

ENV POLL_INTERVAL=2s
ENV SYMBOLS=SBER,GAZP,MOEX
ENV LOG_LEVEL=info

ENTRYPOINT ["/investor"]
