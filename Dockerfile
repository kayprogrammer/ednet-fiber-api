FROM golang:1.24-alpine

RUN mkdir build

# We create folder named build for our app.
WORKDIR /build

COPY ./.env .
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download