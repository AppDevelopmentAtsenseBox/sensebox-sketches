package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type sketchRequestHandler func(http.ResponseWriter, *http.Request) (int, error)

func (fn sketchRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := fn(w, r); err != nil {
		LogInfo("ServeHTTP", "Error: ", err)
		http.Error(w, http.StatusText(status), status)
	}
}

func (sketches *senseBoxSketchesServer) requestHandler(w http.ResponseWriter, req *http.Request) (int, error) {
	LogInfo("requestHandler", "incoming request")
	decoder := json.NewDecoder(req.Body)
	var parsedRequest []SketchRequest

	err := decoder.Decode(&parsedRequest)
	if err != nil {
		LogInfo("requestHandler", "Error decoding JSON payload:", err)
		return http.StatusBadRequest, err
	}

	startContainer()

	return http.StatusOK, nil
}

func (sketches *senseBoxSketchesServer) StartHTTPSServer() {
	LogInfo("StartHTTPSServer", "senseBox Sketches startup")

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(ConfigCaCertBytes); !ok {
		log.Fatalln("Unable to add CA certificate to client certificate pool")
		os.Exit(1)
	}
	LogInfo("StartHTTPSServer", "created client cert pool")

	myServerCertificate, err := tls.X509KeyPair(ConfigServerCertBytes, ConfigServerKeyBytes)
	if err != nil {
		log.Fatalln(err)
	}
	LogInfo("StartHTTPSServer", "imported server cert")

	tlsConfig := &tls.Config{
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                clientCertPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		Certificates:             []tls.Certificate{myServerCertificate},
	}

	tlsConfig.BuildNameToCertificate()
	LogInfo("StartHTTPSServer", "built name to certificate")

	// http.Handle("/", honeybadger.Handler(sketchRequestHandler(sketches.requestHandler)))

	http.Handle("/", sketchRequestHandler(sketches.requestHandler))

	httpServer := &http.Server{
		Addr:      "0.0.0.0:3924",
		TLSConfig: tlsConfig,
	}
	LogInfo("StartHTTPSServer", "configured server")

	LogInfo("StartHTTPSServer", "starting server...")
	log.Fatal(httpServer.ListenAndServeTLS("", ""))
}
