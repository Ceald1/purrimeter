# Use the official Golang image
FROM golang:1.25.1

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./



RUN apt-get update && apt-get install -y ca-certificates

# Download dependencies
RUN go mod download 

# Copy the rest of the application code
COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN echo "updating docs...."

RUN swag init

RUN echo "building API....."

# Build the application
# RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080
ENV REDIS_PASS $REDIS_PASS
ENV SECRET $SECRET
ENV API_ADMIN $API_ADMIN
ENV MANTICORE_HOST $MANTICORE_HOST


# Command to run the application
CMD ["go", "run", "."]

