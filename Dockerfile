# Use the official Golang image as a build stage
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the entire project directory into the container
COPY . .

# Build the Go application
RUN go build -o portfolio

# Use the official Node.js image for building the client
FROM node:18 AS client-builder

# Set the working directory inside the container
WORKDIR /client

# Copy the client directory
COPY client/ .

# Install dependencies and build the client
RUN npm install && npm run build

# Use a minimal base image for the final stage
FROM alpine:latest

# Install necessary CA certificates and other dependencies
RUN apk --no-cache add ca-certificates && apk add --no-cache libc6-compat

# Set the working directory inside the container
WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /app/portfolio .

# Copy the client build output
COPY --from=client-builder /client/dist ./client/dist

# Copy the .env file if needed
COPY .env .env

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./portfolio"]