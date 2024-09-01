# Stage 1: Build the Go application
FROM docker.io/golang:1.22-alpine as build
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY ./go.mod ./go.sum /app/
RUN go mod download

# Copy the rest of the application source code
COPY . /app

# Build the Go application
RUN go build -v -o bin main.go

# Stage 2: Create a minimal image for running the application
FROM docker.io/alpine:3.19
WORKDIR /app

# Install ca-certificates to handle HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the built Go application
COPY --from=build /app/bin /app/bin

# Copy SSL certificates into the container
COPY selfsigned.crt /etc/ssl/certs/selfsigned.crt
COPY selfsigned.key /etc/ssl/private/selfsigned.key

# Expose the port your application listens on
EXPOSE 8080

# Command to run the application
CMD ["/app/bin"]
