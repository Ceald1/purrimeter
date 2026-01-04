# Build stage
FROM golang:1.25.5-trixie AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .
CMD [ "go", "run", "enrichment.go" ]

# Build the binary with optimizations
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main enrichment.go

# # Runtime stage
# FROM alpine:latest

# # RUN apt install -y ca-certificates

# WORKDIR /root/

# # Copy the binary from builder
# COPY --from=builder /app/main .
# RUN chmod +x ./main

# EXPOSE 8080

# CMD ["./main"]