# ---- Build Stage ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (required for some dependencies), and build tools
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go binary
RUN go build -o main .

# ---- Runtime Stage ----
FROM alpine:latest

WORKDIR /root/

# Copy binary from build stage
COPY --from=builder /app/main .

EXPOSE 8000

# Run the binary
CMD ["./main"]