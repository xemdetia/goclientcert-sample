package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
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

// Pulled from golang documentation
func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

// Handler for client cert to check against
func certHandler(w http.ResponseWriter, req *http.Request) {
	if len(req.TLS.VerifiedChains) < 1 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden.\n"))
		return // No chain? Abort
	}
	for _, chain := range req.TLS.VerifiedChains {
		for _, cert := range chain {
			w.Write([]byte(cert.Subject.CommonName + "\n"))
		}
	}
}

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

	clientCaFile, err := ioutil.ReadFile(*tls_ca)
	if err != nil {
		log.Fatal(err)
	}
	stackOfClientCA := x509.NewCertPool()
	stackOfClientCA.AppendCertsFromPEM(clientCaFile)

	tlsContext := &tls.Config{
		Certificates:           []tls.Certificate{cert},
		RootCAs:                stackOfCA,
		ClientCAs:              stackOfClientCA,
		SessionTicketsDisabled: true, // TLS tickets are outside of scope
		MinVersion:             tls.VersionTLS12,
		ClientAuth:             tls.VerifyClientCertIfGiven,
	}

	// Make a copy of the http.Server with the config set
	http.HandleFunc("/", handler)
	http.HandleFunc("/cert", certHandler)
	server := &http.Server{
		Addr:      "127.0.0.1:4443",
		TLSConfig: tlsContext,
	}
	server.ListenAndServeTLS(*tls_certificate, *tls_key)
}
