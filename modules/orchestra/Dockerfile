# Use the official Golang image
FROM golang:1.25.1

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Update and install certs
RUN apt-get update && apt-get install -y ca-certificates

# Download dependencies
RUN go mod download 

# Copy the rest of the application code
COPY . .
ENV CONDUCTOR_HOST $CONDUCTOR_HOST
ENV CONDUCTOR_PORT $CONDUCTOR_PORT
ENV MANTICORE_HOST $MANTICORE_HOST
ENV MANTICORE_PORT $MANTICORE_PORT

CMD ["go", "run", "orchestra.go"]