# --------------------------------------------------------
# Stage 1: The Builder
# Use the official Golang image to compile the source code
# --------------------------------------------------------
FROM golang:1.25-alpine AS builder

# Install git (required for fetching some Go dependencies)
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first to leverage Docker caching
# If go.mod hasn't changed, Docker skips 'go mod download'
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
# -o main: Name the output binary "main"
# ./cmd/api: Point this to where your main.go lives
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/Server

# --------------------------------------------------------
# Stage 2: The Runtime
# Use a tiny Alpine Linux image for production
# --------------------------------------------------------
FROM alpine:latest

# Install certificates so the app can make HTTPS calls if needed
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the compiled binary from the Builder stage
COPY --from=builder /app/main .

# Expose the port your app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]