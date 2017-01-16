package main

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	// "encoding/json"
	"log"
	"os"

	"fmt"

	"github.com/fsouza/go-dockerclient"
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

	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}

	pullImageErr := client.PullImage(docker.PullImageOptions{Repository: "mpfeil/docker-frab"}, docker.AuthConfiguration{})
	if pullImageErr != nil {
		panic(pullImageErr)
	}
	fmt.Println("Image pulled!")

	// exposedCadvPort := map[docker.Port]struct{}{"3000/tcp": {}}

	createContConf := docker.Config{
		// ExposedPorts: exposedCadvPort,
		Image: "mpfeil/docker-frab",
	}

	portBindings := map[docker.Port][]docker.PortBinding{
		"3000/tcp": {{HostIP: "0.0.0.0", HostPort: "3000"}},
	}

	createContHostConfig := docker.HostConfig{
		// Binds:           []string{"/var/run:/var/run:rw", "/sys:/sys:ro", "/var/lib/docker:/var/lib/docker:ro"},
		PortBindings:    portBindings,
		PublishAllPorts: false,
		Privileged:      false,
	}

	createContOps := docker.CreateContainerOptions{
		Name:       "frab",
		Config:     &createContConf,
		HostConfig: &createContHostConfig,
	}

	// fmt.Printf("\nops = %s\n", createContOps)

	cont, err := client.CreateContainer(createContOps)
	if err != nil {
		fmt.Printf("create error = %s\n", err)
	}
	// fmt.Printf("container = %s\n", cont)

	err = client.StartContainer(cont.ID, nil)
	if err != nil {
		fmt.Printf("start error = %s\n", err)
	}
	fmt.Printf("start = %s\n", err)

	// container, createContainerErr := client.CreateContainer(docker.CreateContainerOptions{Name: "docker-frab", Config: &docker.Config{Image: "mpfeil/docker-frab"}})
	// if createContainerErr != nil {
	// 	panic(createContainerErr)
	// }
	// fmt.Println("Container created: " + container.ID)

	// portBindings := []docker.PortBinding{}
	// portBindings = append(portBindings, docker.PortBinding{HostIP: "", HostPort: "3000"})

	// m := make(map[docker.Port][]docker.PortBinding)
	// m["3000/tcp"] = portBindings

	// // fmt.Println(m)

	// hostConfig := docker.HostConfig{
	// 	PortBindings: m,
	// }
	// fmt.Println(hostConfig)

	// startContainerErr := client.StartContainer(container.ID, &hostConfig)
	// if startContainerErr != nil {
	// 	panic(startContainerErr)
	// }
	// fmt.Println("Container started")

	// reader, imagePullErr := cli.ImagePull(context.Background(), "abiosoft/caddy", types.ImagePullOptions{})
	// if imagePullErr != nil {
	// 	panic(imagePullErr)
	// }

	// _, err = io.Copy(os.Stdout, reader)
	// if err != nil {
	// 	panic(err)
	// }
	// reader.Close()

	// _, containerCreateErr := cli.ContainerCreate(context.Background(), &container.Config{Image: "abiosoft/caddy"}, nil, nil, "sketches")
	// if containerCreateErr != nil {
	// 	panic(containerCreateErr)
	// }

	// containerStartErr := cli.ContainerStart(context.Background(), "sketches", types.ContainerStartOptions{})
	// if containerStartErr != nil {
	// 	panic(containerStartErr)
	// }
	// fmt.Println("Container started")

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
