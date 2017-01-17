package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsouza/go-dockerclient"
)

var client *docker.Client
var container *docker.Container
var hostConfig *docker.HostConfig

/*
{
  "payload": {
    "box": {
      "id": "<some valid senseBox id>",
      "sensors": [
        {
          "title": "<some title>",
          "type": "<some type>",
          "id": "<some valid senseBox sensor id>"
        },
        ...
      ]
    }
  },
}
*/

type SketchRequest struct {
	NetworkType string                 `json:"networktype"`
	Payload     map[string]interface{} `json:"payload"`
}

func initContainer() {
	endpoint := "unix:///var/run/docker.sock"
	cli, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	LogInfo("Sketches - InitContainer", "Client created")
	client = cli

	pullImageErr := client.PullImage(docker.PullImageOptions{Repository: "mpfeil/docker-frab"}, docker.AuthConfiguration{})
	if pullImageErr != nil {
		panic(pullImageErr)
	}
	LogInfo("Sketches - InitContainer", "Image pulled")

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
	hostConfig = &createContHostConfig

	createContOps := docker.CreateContainerOptions{
		Name:       "frab",
		Config:     &createContConf,
		HostConfig: &createContHostConfig,
	}

	cont, createContainerErr := client.CreateContainer(createContOps)
	if createContainerErr != nil {
		log.Fatalln("Container was not created: ")
		os.Exit(1)
	}

	container = cont
}

func startContainer() {
	startContainerErr := client.StartContainer(container.ID, hostConfig)
	if startContainerErr != nil {
		fmt.Printf("start error = %s\n", startContainerErr)
	}
	fmt.Printf("start = %s\n", startContainerErr)
}
