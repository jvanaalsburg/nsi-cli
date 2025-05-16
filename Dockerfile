FROM golang:1.23

WORKDIR /app

# Install Go dependencies.
COPY src/go.mod src/go.sum .
RUN go mod download
