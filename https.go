package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	// "encoding/json"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	reader, imagePullErr := cli.ImagePull(context.Background(), "abiosoft/caddy", types.ImagePullOptions{})
	if imagePullErr != nil {
		panic(imagePullErr)
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		panic(err)
	}
	reader.Close()

	_, containerCreateErr := cli.ContainerCreate(context.Background(), &container.Config{Image: "abiosoft/caddy"}, nil, nil, "sketches")
	if containerCreateErr != nil {
		panic(containerCreateErr)
	}

	containerStartErr := cli.ContainerStart(context.Background(), "sketches", types.ContainerStartOptions{})
	if containerStartErr != nil {
		panic(containerStartErr)
	}
	fmt.Println("Container started")

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
