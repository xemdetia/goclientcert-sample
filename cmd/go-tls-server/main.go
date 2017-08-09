package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
)

var (
	tls_certificate = flag.String("cert",
		"host.crt",
		"The HTTPS server's certificate.")
	tls_key = flag.String("key",
		"host.key",
		"The HTTPS server's private key.")
	tls_ca = flag.String("CAfile",
		"ca.crt",
		"The stacked PEM representing the -cert certificate chain.")
	tls_client_ca = flag.String("clientCAfile",
		"ca.crt",
		"The stacked PEM representing possible clients. If not found uses the -CAfile instead.")
	subject_regex = flag.String("e",
		".*",
		"Regex expression passed in to match the subject or subjectAltName of a certificate. By default it matches everything.")
)

func main() {
	flag.Parse()
	cert, err := tls.LoadX509KeyPair(*tls_certificate, *tls_key)
	if err != nil {
		log.Fatal(err)
	}

	caFile, err := ioutil.ReadFile(*tls_ca)
	if err != nil {
		log.Fatal(err)
	}
	stackOfCA := x509.NewCertPool()
	stackOfCA.AppendCertsFromPEM(caFile)

	tlsContext := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      stackOfCA,
	}
	tlsContext.BuildNameToCertificate()
}
