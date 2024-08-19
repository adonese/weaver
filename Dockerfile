# Use the official Go image as a base image for building
FROM golang:1.23-alpine as build

# Set the Current Working Directory inside the container
WORKDIR /app

# Install necessary dependencies
RUN apk add --no-cache git

# Install weaver
RUN go install github.com/ServiceWeaver/weaver/cmd/weaver@latest

# Copy the go.mod and go.sum files first to enable dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o payment_system .

# Use a minimal image to run the application
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=build /app/payment_system .
COPY --from=build /app/weaver.toml .
COPY --from=build /app/redoc-static.html .
COPY --from=build /go/bin/weaver /usr/local/bin/weaver

# Expose the port that the application will run on
EXPOSE 8080

# Command to run the executable
CMD ["/usr/local/bin/weaver", "single", "deploy", "weaver.toml"]
