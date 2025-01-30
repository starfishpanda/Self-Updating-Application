FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build


# Build Server
WORKDIR /build/server
COPY server/go.mod server/go.sum* ./
COPY server/server.go .
RUN CGO_ENABLED=0 go build -o server

# Build current and update version of client
WORKDIR /build/client
COPY client/go.mod client/go.sum* ./
COPY client/client.go .

#Build version 1.1.1
RUN CGO_ENABLED=0 go build -ldflags "-X main.currentVersion=1.1.1" -o client-1.1.1
#Build update version 1.1.2
RUN CGO_ENABLED=0 go build -ldflags "-X main.currentVersion=1.1.2" -o client-1.1.2

FROM alpine:latest

# In case we need ca certs for HTTPS
# RUN apk add --no-cache ca-certificates

# Create directory structure in container
WORKDIR /app
RUN mkdir -p server/binaries client && \
  chown -R nobody:nobody /app && \
  chmod -R 755 /app



# Copy server binary from build
COPY --from=builder /build/server/server ./server/

# Copy client and update binaries from build
COPY --from=builder /build/client/client-1.1.1 ./client/client
COPY --from=builder /build/client/client-1.1.2 ./server/binaries/myapp-update

# Make start script and binaries executable
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh && \
  chmod +x /app/server/server && \
  chmod +x /app/client/client && \
  chmod +x /app/server/binaries/myapp-update

USER nobody

WORKDIR /app
CMD ["./start.sh"]