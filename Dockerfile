# Use a lightweight alpine image as the base
FROM alpine:latest

# Install any required dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory for the container
WORKDIR /app

# Copy the main.go file to the container
COPY /api/main.go .

# Build the Go binary
RUN apk --no-cache add go && \
    go build -o main .

# Set the entrypoint for the container
ENTRYPOINT ["./main"]
