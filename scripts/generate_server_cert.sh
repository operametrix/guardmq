#!/bin/bash

# Remove previous certificate
rm certs/server*

# Generate the keys for the server
openssl genrsa -out certs/server.key 4096

# Generate the CSR for the server
openssl req -new -key certs/server.key -out certs/server.csr

# Sign the server certificate with the CA
openssl x509 -req -in certs/server.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/server.crt
