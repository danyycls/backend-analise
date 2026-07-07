FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/inicializador_banco ./cmd/migrate/

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/inicializador_banco /app/inicializador_banco
COPY internal/shared/migrations /app/internal/shared/migrations
COPY .env /app/.env

CMD ["./inicializador_banco"]
