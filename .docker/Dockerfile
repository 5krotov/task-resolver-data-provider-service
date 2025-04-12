# --- BUILD STAGE ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o data-provider-service ./cmd/main.go

# --- RUN STAGE ---
FROM alpine:latest AS release

COPY --from=builder /app/data-provider-service /app/data-provider-service/
COPY config.yaml /etc/data-provider-service/

WORKDIR /app/data-provider-service/

EXPOSE 8080

CMD ["./data-provider-service"]