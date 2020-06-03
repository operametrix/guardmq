#!/bin/bash

# Generate the keys for the CA
openssl genrsa -out ca.key 4096

# Generate the certificate of the CA
openssl req -new -x509 -key ca.key -out ca.crt

# Generate the keys for the server
openssl genrsa -out server.key 4096

# Generate the CSR for the server
openssl req -new -key server.key -out server.csr

# Sign the server certificate with the CA
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt
