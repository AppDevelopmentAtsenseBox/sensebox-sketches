package main

import (
	"fmt"
	"log"
	"os"

	"strings"

	"github.com/fsouza/go-dockerclient"
)

var client *docker.Client
var containerID string
var hostConfig *docker.HostConfig

const repository = "mpfeil/docker-frab"

/* Payload that is send by the request.
{
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
}
*/
type SketchRequest struct {
	Box Box `json:"box"`
}

/**/
type Box struct {
	ID      string   `json:"_id"`
	Sensors []Sensor `json:"sensors"`
	Model   string   `json:"model"`
}

/**/
type Sensor struct {
	Title      string `json:"title"`
	Type       string `json:"type"`
	ID         string `json:"_id"`
	SensorType string `json:"sensorType"`
}

func initContainer() {
	endpoint := "unix:///var/run/docker.sock"
	cli, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	LogInfo("Sketches - InitContainer", "Client created")
	client = cli

	// Check Images
	images, listImagesErr := client.ListImages(docker.ListImagesOptions{})
	if listImagesErr != nil {
		panic(listImagesErr)
	}

	var imageExists bool
	for _, element := range images {
		for _, elem := range element.RepoTags {
			if strings.HasPrefix(elem, repository) {
				imageExists = true
				LogInfo("Sketches - InitContainer", "Image already exists!")
			}
		}
	}

	if !imageExists {
		LogInfo("Sketches - InitContainer", "Image not exisitng. Imgage pulling...")
		pullImageErr := client.PullImage(docker.PullImageOptions{Repository: "mpfeil/docker-frab"}, docker.AuthConfiguration{})
		if pullImageErr != nil {
			panic(pullImageErr)
		}
		LogInfo("Sketches - InitContainer", "Image pulled")
	}

	// Check Containers
	containers, listContainersError := client.ListContainers(docker.ListContainersOptions{All: true})
	if listContainersError != nil {
		panic(listContainersError)
	}

	var containerExists bool
	for _, container := range containers {
		if container.Image == repository {
			containerExists = true
			containerID = container.ID
			LogInfo("Sketches - InitContainer", "Container already exists!")
		}
	}

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

	if !containerExists {
		LogInfo("Sketches - InitContainer", "Container does not exists. New Container is creating...")

		cont, createContainerErr := client.CreateContainer(createContOps)
		if createContainerErr != nil {
			log.Fatalln("Container was not created: ")
			os.Exit(1)
		}
		LogInfo("Sketches - InitContainer", "Container is created!")

		containerID = cont.ID
	}
}

func startContainer() {
	startContainerErr := client.StartContainer(containerID, hostConfig)
	if startContainerErr != nil {
		fmt.Printf("start error = %s\n", startContainerErr)
	}
	LogInfo("Sketches - StartContainer", "Container started!")
}

func stopContainer() {
	stopContainerErr := client.StopContainer(containerID, 5000)
	if stopContainerErr != nil {
		panic(stopContainerErr)
	}
	LogInfo("Sketches - StopContainer", "Container stopped!")
}
