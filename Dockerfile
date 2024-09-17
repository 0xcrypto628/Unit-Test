# Dockerfile

# Step 1: Use Go image for build
FROM golang:1.23 AS builder

# Set environment variables
ENV GO111MODULE=on
WORKDIR /app

# Copy source code to the build container
COPY . .

# Build the blockchain binary
RUN go mod tidy
RUN go build -o mychaind ./cmd/mychaind

# Step 2: Create a lightweight image for running the binary
FROM alpine:3.18

WORKDIR /root/

# Copy the binary from the builder container
COPY --from=builder /app/mychaind /usr/local/bin/mychaind

# Expose ports for communication (P2P, RPC, API)
EXPOSE 26656 26657 1317

# Command to start the blockchain node
CMD ["mychaind", "start"]
