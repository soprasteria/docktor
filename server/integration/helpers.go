package integration

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/moby/moby/client"
	"github.com/moby/moby/container"
	log "github.com/sirupsen/logrus"
	"github.com/soprasteria/docktor/cmd"
)

const (
	dockerHost = "unix:///var/run/docker.sock"
	dockerAPIVersion = "1.23"
	tcpConnection = "tcp"
	allIP = "0.0.0.0"
)

// InstantiateDocker instantiates a daemon docker
func InstantiateDocker() (*client.Client, error) {
	docker, err := client.NewClient(dockerHost, dockerAPIVersion, nil, make(map[string]string))
	if err != nil {
		return nil, err
	}
	return docker, nil
}

// PullImage pulls a docker image
func PullImage(docker *client.Client, imageName, imageVersion string) error {

	if docker == nil {
		return fmt.Errorf("docker daemon is nil, the image can't be pulled")
	}
	if imageName == "" {
		return fmt.Errorf("image's name can't be empty")
	}
	if imageVersion == "" {
		return fmt.Errorf("image's version can't be empty")
	}
	_, err := docker.ImagePull(
		context.Background(),
		fmt.Sprintf("%v:%v", imageName, imageVersion),
		types.ImagePullOptions {}
	)
	return  err
}

// CreateContainer creates a container from an docker image
func CreateContainer(docker *client.Client, containerName, portBinding, imageName, imageVersion string) (container.ContainerCreateCreatedBody, error) {

	if docker == nil {
		return nil, fmt.Errorf("docker daemon is nil, the container can't be created")
	}
	portBindings := nat.PortMap{}
	port := fmt.Sprintf("%v/%s", portBinding, tcpConnection)
	portBindings[port] = []nat.PortBinding{{HostIP: allIP, HostPort: portBinding}}

	container, err := docker.ContainerCreate(
		context.Background(),
		containertypes.Config {
			Image: fmt.Sprintf("%v:%v", imageName, imageVersion)
		},
		containertypes.HostConfig {
			PortBindings: portBindings
		},
		network.NetworkingConfig {},
		containerName)
	if err != nil {
		return nil, err
	}
	return container, nil
}

// StartContainer starts a container
func StartContainer(docker *client.Client, container *container.ContainerCreateCreatedBody) error {

	if docker == nil {
		return nil, fmt.Errorf("docker daemon is nil, the container can't be started")
	}
	if container == nil {
		return nil, fmt.Errorf("container is nil, the container can't be started")
	}
	err := docker.ContainerStart(
		context.Background(),
		container.ID,
		types.ContainerStartOptions {}
	)
	return err
}

/*
// InstantiateDocker instantiates a daemon docker
func InstantiateDocker() (*dockerapi.Client, error) {

	log.Info("Instantiate docker daemon")
	// Instantiate docker daemon
	docker, err := dockerapi.NewClient(dockerHost)
	if err != nil {
		return nil, err
	}
	return docker, nil
}*/

// RunContainer runs a container (pull image + create container + start container)
func RunContainer(docker *client.Client, containerName, imageName, imageVersion, port string) (*container.Container, error) {

	if docker == nil {
		return nil, fmt.Errorf("docker daemon is nil, the container can't be instantiated")
	}

	/*docker.ContainerCreate(
	context.Background(),
	&container.Config{
		Image: fmt.Sprintf("%s%s%s", imageName, ":", imageVersion)
	},
	&container.HostConfig{
		PortBindings: ["12345"]
	},
	nil,
	name)*/

	// Start container
	container.
}

/*
// InstantiateContainer instantiates a container with its name, image name, image version and port used
func InstantiateContainer(docker *dockerapi.Client, name, imageName, imageVersion, port string) (*dockerapi.Container, error) {

	if docker == nil {
		return nil, fmt.Errorf("docker daemon is nil, the container can't be instantiated")
	}

	portBinding := dockerapi.PortBinding{HostPort: port, ContainerPort: port}

	container, err := docker.NewContainer(dockerapi.ContainerOptions{
		Image:        fmt.Sprintf("%s%s%s", imageName, ":", imageVersion),
		Name:         name,
		PortBindings: []dockerapi.PortBinding{portBinding},
	})

	if err != nil {
		return nil, err
	}
	return container, nil
}*/
/*
// RunContainer runs a container
func RunContainer(container *dockerapi.Container) error {
	return container.Run(false)
}*/

// StartServer starts the docktor server
func StartServer() {

	log.Info("Start docktor server")
	// Run the docktor server
	go func() {
		cmd.ServeCmd.Run(nil, []string{})
	}()

	// Wait until server return Status code OK
	serverStarted := false
	count := 0
	for !serverStarted {
		time.Sleep(100 * time.Millisecond)
		count += 100
		log.Info("Checking if docktor server is started...")
		resp, err := http.Get("http://localhost:8080/ping")
		if err != nil {
			log.Info("Server not started:", err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Info("Status Not OK:", resp.StatusCode)
			continue
		}
		serverStarted = true
	}
	log.Info(fmt.Sprintf("%s%d%s", "Docktor server STARTED in ", count, " ms"))
}

/*
// CleanContainer cleans a container by stopping it and removing it
func CleanContainer(container *dockerapi.Container) error {

	if container == nil {
		return fmt.Errorf("container is nil, it can't be cleaned")
	}
	log.Info(fmt.Sprintf("Cleaning container ID : %v", container.Container.ID))
	err := container.StopAndRemove(true)
	if err != nil {
		return err
	}
	return nil
}*/
