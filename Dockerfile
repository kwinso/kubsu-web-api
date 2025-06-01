# Stage 1: Build the Go application
FROM golang:1.24 AS builder

# Set the Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app (adjust -o to your desired binary name)
RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go

# Stage 2: Create a minimal Docker image
FROM gcr.io/distroless/static-debian12

# Copy the built binary from the builder stage
COPY --from=builder /server /server

# Expose the port your app runs on
EXPOSE 8080

# Command to run the binary when the container starts
ENTRYPOINT ["/server"]