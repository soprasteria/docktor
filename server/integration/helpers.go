package integration

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/soprasteria/dockerapi"
	"github.com/soprasteria/docktor/cmd"
)

const (
	dockerHost = "unix:///var/run/docker.sock"
)

// InstantiateDocker instantiates a daemon docker
func InstantiateDocker() (*dockerapi.Client, error) {

	log.Info("Instantiate docker daemon")
	// Instantiate docker daemon
	docker, err := dockerapi.NewClient(dockerHost)
	if err != nil {
		return nil, err
	}
	return docker, nil
}

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
}

// RunContainer runs a container
func RunContainer(container *dockerapi.Container) error {
	return container.Run(false)
}

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

// CleanContainer cleans a container by stopping it and removing it
func CleanContainer(container *dockerapi.Container) error {

	if container == nil {
		return fmt.Errorf("container is nil, it can't be cleaned")
	}
	fmt.Println("test1")
	err := container.StopAndRemove(true)
	if err != nil {
		return err
	}
	fmt.Println("test2")
	return nil
}
