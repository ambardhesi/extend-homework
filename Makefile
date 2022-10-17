.PHONY: binaries clean vendor ca_key ca_cert svr_key svr_csr svr_cert ambar_key ambar_csr ambar_cert \
	bob_key bob_csr bob_cert

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
ROOT=$(shell pwd)
COMMANDS=$(shell go list ./... | grep cmd/)

OPENSSL?=/usr/bin/openssl
CERT_DIR=$(ROOT)/certs
CA_CERT="$(CERT_DIR)/ca-cert.pem"
CA_KEY="$(CERT_DIR)/ca-key.pem"
SVR_CERT="$(CERT_DIR)/svr-cert.pem"
SVR_KEY="$(CERT_DIR)/svr-key.pem"
SVR_CSR="$(CERT_DIR)/svr-csr.csr"
AMBAR_CERT="$(CERT_DIR)/ambar-cert.pem"
AMBAR_KEY="$(CERT_DIR)/ambar-key.pem"
AMBAR_CSR="$(CERT_DIR)/ambar-csr.csr"
BOB_CERT="$(CERT_DIR)/bob-cert.pem"
BOB_KEY="$(CERT_DIR)/bob-key.pem"
BOB_CSR="$(CERT_DIR)/bob-csr.csr"

binaries:
	$(GOBUILD) -o "." $(COMMANDS) 	

clean:
	rm -rf vendor
	rm certs/*.pem
	rm certs/*.csr
	rm certs/*.srl
	$(GOCLEAN)

vendor:
	$(GOCMD) mod vendor

certs: ca_cert svr_cert ambar_cert bob_cert 

ca_key: 
	$(OPENSSL) genpkey -algorithm ed25519 -out $(CA_KEY) -outform PEM

ca_cert: ca_key
	$(OPENSSL) req -x509 -newkey rsa:4096 -key $(CA_KEY) -out $(CA_CERT) -subj "/C=CA/ST=BC/L=Vancouver/OU=SignedCA/CN=localhost/emailAddress=foo@foo.com"

svr_key:
	$(OPENSSL) genpkey -algorithm ed25519 -out $(SVR_KEY) -outform PEM

svr_csr: svr_key
	$(OPENSSL) req -new -key $(SVR_KEY) -out $(SVR_CSR) -config "$(CERT_DIR)/svr.conf" 

svr_cert: svr_csr
	$(OPENSSL) x509 -req -in $(SVR_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -out $(SVR_CERT) -extfile $(CERT_DIR)/svr.conf -extensions my_extensions

ambar_key:
	$(OPENSSL) genpkey -algorithm ed25519 -out $(AMBAR_KEY) -outform PEM

ambar_csr: ambar_key
	$(OPENSSL) req -new -key $(AMBAR_KEY) -out $(AMBAR_CSR) -subj "/CN=ambar.dhesi@gmail.com"

ambar_cert: ambar_csr
	$(OPENSSL) x509 -req -in $(AMBAR_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -out $(AMBAR_CERT) 

bob_key: 
	$(OPENSSL) genpkey -algorithm ed25519 -out $(BOB_KEY) -outform PEM

bob_csr: bob_key
	$(OPENSSL) req -new -key $(BOB_KEY) -out $(BOB_CSR) -subj "/CN=bob"

bob_cert: bob_csr
	$(OPENSSL) x509 -req -in $(BOB_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -out $(BOB_CERT) 
