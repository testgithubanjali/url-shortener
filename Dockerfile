# -----------------------------
# Build Stage
# -----------------------------
FROM golang:1.26.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd

# -----------------------------
# Runtime Stage
# -----------------------------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/url-shortener .

COPY .env .

EXPOSE 8081

CMD ["./url-shortener"]