# syntax=docker/dockerfile:1

# Build Stage
FROM golang:1.21 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Generate Ent code
RUN go generate ./ent

# Build the binary
RUN go build -o main .

# Final Stage
FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=builder /app/main .

CMD ["/app/main"]