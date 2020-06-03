#!/bin/bash

mkdir -p certs/peers

PEER="test"

# Generate the keys for the peer
openssl genrsa -out certs/peers/$PEER.key 4096

# Generate the CSR for the peer
openssl req -new -key certs/peers/$PEER.key -out certs/peers/$PEER.csr

# Sign the peer certificate with the CA
openssl x509 -req -in certs/peers/$PEER.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/peers/$PEER.crt
