#!/bin/bash

# Generate the keys for the CA
openssl genrsa -out certs/ca.key 4096

# Generate the certificate of the CA
openssl req -new -x509 -key certs/ca.key -out certs/ca.crt