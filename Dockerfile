# Step 1: Use an official Go image for building the binary
FROM golang:1.24 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first (for caching deps)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build the Go app (binary called `app`)
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Step 2: Use a smaller image for running (distroless or alpine)
FROM alpine:latest

# Add a non-root user for security
RUN adduser -D -g '' appuser

WORKDIR /app
COPY --from=builder /app/app .

# Run the binary as non-root
USER appuser

# Expose your app port (change if you use a different port)
EXPOSE 8080

# Start the app
ENTRYPOINT ["./app"]
