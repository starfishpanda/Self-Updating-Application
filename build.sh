#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Define binaries version variables
INITIAL_VERSION="1.1.1"
UPDATED_VERSION="1.1.2"

# Build the initial client binary in client directory
echo "Building initial client version $INITIAL_VERSION..."
cd client
go build -ldflags "-X main.currentVersion=$INITIAL_VERSION" -o myapp client.go
cd ..

# Build the updated client binary in binaries directory
echo "Building updated client version $UPDATED_VERSION..."
cd client
go build -ldflags "-X main.currentVersion=$UPDATED_VERSION" -o ../server/binaries/myapp-update client.go
cd ..

# Build the server binary
echo "Building server..."
cd server
go build -o server server.go
cd ..

echo "Build process completed successfully."
