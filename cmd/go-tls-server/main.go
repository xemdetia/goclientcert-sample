package main

import "flag"

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
}
