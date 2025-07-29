# syntax=docker/dockerfile:1

# Stage 1: Build
FROM golang:1.24.1-alpine AS builder

# Needed for static linking
RUN apk add --no-cache upx

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Regenerate ent if needed
RUN go generate ./ent

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Optional: compress binary (smaller image)
RUN upx --best main

# Stage 2: Minimal final image
FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

# Run the binary
ENTRYPOINT ["/app/main"]