# Use the official Golang image to create a build artifact
FROM golang:alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install necessary packages including tzdata for time zone data
RUN apk update && apk add --no-cache tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o argus .

# Start a new stage from scratch
FROM alpine:latest

# Install tzdata package for time zone data
RUN apk update && apk add --no-cache tzdata

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/argus .

# Copy the config file
COPY config/config.json /root/config/config.json

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./argus"]