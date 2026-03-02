# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/promptkit ./cmd/promptkit/

# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/promptkit /usr/local/bin/promptkit
COPY --from=builder /app/templates/ ./templates/

ENTRYPOINT ["promptkit"]
CMD ["--help"]
