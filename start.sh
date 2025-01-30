#!/bin/sh

# Start the update server in the background
cd /app/server
./server &

# Start the initial client
cd /app/client
./client
