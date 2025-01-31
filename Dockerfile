FROM golang:1.20-alpine AS builder

RUN go version

# Install any needed packages
RUN apk add --no-cache \
  build-base \
  git \
  ca-certificates \
  gcc \
  musl-dev

ENV CGO_ENABLED=1
ENV GO111MODULE=on

WORKDIR /build
RUN go mod init build

# Build the Server and run Keygen
WORKDIR /build/server
COPY server/go.mod ./
RUN go mod download
COPY server/ ./



WORKDIR /build/server/keygen
RUN go run keygen.go
WORKDIR /build/server


RUN go build -o server .

# Build the Client
WORKDIR /build/client
COPY client/ ./
RUN go list -m all
RUN go env && ls -la /usr/local/go/src/crypto

# 3) Build two client versions
RUN go build -v -ldflags "-X main.currentVersion=1.1.1" -o client-1.1.1
RUN go build -v -ldflags "-X main.currentVersion=1.1.2" -o client-1.1.2

FROM alpine:latest

# Create directory structure in container
WORKDIR /app
RUN mkdir -p server/binaries client \
  && chown -R nobody:nobody /app \
  && chmod -R 755 /app

# Copy server binary
COPY --from=builder /build/server/server ./server/

# Copy PEM files generated by keygen (private key for server, public key for client)
COPY --from=builder /build/server/keygen/private.pem ./server/
COPY --from=builder /build/server/keygen/public.pem  ./client/

# Copy client + update binaries
COPY --from=builder /build/client/client-1.1.1 ./client/client
COPY --from=builder /build/client/client-1.1.2 ./server/binaries/myapp-update

# Copy startup script, etc.
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh \
  && chmod +x /app/server/server \
  && chmod +x /app/client/client \
  && chmod +x /app/server/binaries/myapp-update \
  # Set restrictive permissions on private key
  && chmod 600 /app/server/private.pem \
  # Public key can be world-readable
  && chmod 644 /app/client/public.pem \
  # Ensure 'nobody' owns everything
  && chown -R nobody:nobody /app

USER nobody
CMD ["/app/start.sh"]
