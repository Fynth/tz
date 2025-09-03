# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

RUN wget -qO- https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | \
    tar -xzf - -C /bin migrate

# Финальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/



COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /bin/migrate /bin/migrate
COPY docs ./docs

EXPOSE 8080
CMD ["./server"]
