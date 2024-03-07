FROM golang:alpine AS builder
# specify the name of working dir
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY ../. .

# Build the Go application
RUN go build -o urlshortner

EXPOSE 1323

# Command to run your application
CMD ["./urlshortner"]