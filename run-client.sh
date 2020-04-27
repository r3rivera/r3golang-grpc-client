#!/bin/bash
echo "Running the Golang GRPC Server"


export SECURE_MODE=true
export CLIENT_SERVER_PORT=50051
export CLIENT_SERVER_HOST=localhost
export CLIENT_CERT_PATH="./certs/r2play.ca.crt"

echo "Client Server Port Connection is ${CLIENT_SERVER_PORT}"
echo "Client Server Host Name is ${CLIENT_SERVER_HOST}"
echo "Client Certificate Path is ${CLIENT_CERT_PATH}"
go run basic-client.go