# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/discord ./cmd

# Runtime stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/discord /usr/local/bin/discord
COPY --from=builder /app/migrations /app/migrations
WORKDIR /app

EXPOSE 8080
ENV PORT=8080
CMD ["/usr/local/bin/discord"]
