FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o articleforum ./cmd/server

FROM alpine:latest

RUN apk add --no-cache ca-certificates postgresql-client

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/articleforum .

COPY --from=builder /app/migrations ./migrations

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./articleforum", "-storage", "postgres"]