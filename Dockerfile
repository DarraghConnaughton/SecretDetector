# Builder Stage
FROM golang:latest AS builder

# Set the working directory in the builder image
WORKDIR /secretdetector

# Copy the Go source code and Makefile
COPY . .

# Build the Golang binary
RUN make build

# Final Stage
FROM ubuntu:latest

# Set the working directory in the final image
WORKDIR /secretdetector

# Copy the binary from the builder image to the final image
COPY --from=builder /secretdetector/releases/* /secretdetector/

# Example: expose a port if your application listens on a specific port
# EXPOSE 8080

# Run the binary as the default command
CMD ["./secretdetector"]