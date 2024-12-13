# Use an official Go runtime as a parent image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source files
COPY . .

# Build the Go app
RUN go build -o p2p_network_service .

# Expose the port the app runs on
EXPOSE 10000

# Run the binary program produced by `go build`
CMD ["./p2p_network_service"]
