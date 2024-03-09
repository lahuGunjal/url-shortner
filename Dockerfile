# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Specify the name of the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o urlshortener

# Stage 2: Production image
FROM alpine:latest

# Specify the name of the working directory
WORKDIR /app

# Copy only the necessary artifacts from the previous stage
COPY --from=builder /app/urlshortener .

# Expose the port the application runs on
EXPOSE 1323

# Command to run your application
CMD ["./urlshortener"]
