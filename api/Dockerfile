# Use the official Golang image
FROM golang:1.25.1

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on

RUN apt-get update && apt-get install -y ca-certificates

# Download dependencies
RUN go mod download 

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8000

# Command to run the application
CMD ["./main"]

