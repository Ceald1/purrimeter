# Build stage
FROM golang:1.25.5-trixie AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

CMD ["go", "run", "alerts.go"]
