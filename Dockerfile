# Builder Stage
FROM golang:latest AS builder

# Set the working directory in the builder image
WORKDIR /cmd

# Copy the Go source code and Makefile
COPY . .
RUN make build

# Move executable to final image.
FROM ubuntu:latest

# Starting point of the recursive scan.
ENV SCAN_START_DIRECTORY "thirdParty"

# Set the working directory in the final image
WORKDIR /cmd
# Copy the binary from the builder image to the final image
COPY --from=builder /cmd/releases/secretdetector /cmd/
COPY --from=builder /cmd/data/ /cmd/data/
COPY --from=builder /cmd/$SCAN_START_DIRECTORY/ /cmd/$SCAN_START_DIRECTORY/

# Create a non-root user and set permissions
RUN groupadd -r secretdetector && useradd -r -g secretdetector secretdetector
RUN chown -R secretdetector:secretdetector /cmd
USER secretdetector
CMD ["/cmd/secretdetector"]
