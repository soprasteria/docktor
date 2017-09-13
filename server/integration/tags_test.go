package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/soprasteria/dockerapi"
	"github.com/soprasteria/docktor/cmd"
)

const (
	dockerHost  = "unix:///var/run/docker.sock"
	docktorHost = "mongodb://localhost:%v/docktor"
)

var (
	docktorMongo *dockerapi.Container
)

func TestServer(t *testing.T) {

	log.Info("Instanciate docktor mongoDB")
	// Instanciate mongoDB
	docker, err := dockerapi.NewClient(dockerHost)
	Convey("When a new client is created", t, func() {
		Convey("Then no error shoud be returned", func() {
			So(err, ShouldBeNil)
		})
		Convey("The client docker should not be nil", func() {
			So(docker, ShouldNotBeNil)
		})
	})
	portBinding := dockerapi.PortBinding{HostPort: "27017", ContainerPort: "27017"}

	c, err := docker.NewContainer(dockerapi.ContainerOptions{
		Image:        "mongo:latest",
		Name:         "mongo-test",
		PortBindings: []dockerapi.PortBinding{portBinding},
	})
	Convey("When a new container for mongoDB is created", t, func() {
		Convey("Then no error shoud be returned", func() {
			So(err, ShouldBeNil)
		})
		Convey("The container should not be nil", func() {
			So(c, ShouldNotBeNil)
		})
	})

	log.Info("Run docktor mongoDB")
	// Run docktor mongoDB
	docktorMongo = c
	err = docktorMongo.Run(false)
	Convey("When the container for mongoDB is running", t, func() {
		Convey("Then no error shoud be returned", func() {
			So(err, ShouldBeNil)
		})
	})

	log.Info("Run docktor server")
	// Run the docktor server
	go func() {
		cmd.ServeCmd.Run(nil, []string{})
	}()

	// Wait until server return Status code OK
	serverStarted := false
	count := 0
	for !serverStarted {
		time.Sleep(10 * time.Millisecond)
		count += 10
		log.Println("Checking if docktor server is started...")
		resp, err := http.Get("http://localhost:8080/ping")
		if err != nil {
			log.Println("Server not started:", err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Println("Status Not OK:", resp.StatusCode)
			continue
		}
		serverStarted = true
	}
	log.Info(fmt.Sprintf("%s%d%s", "Docktor server STARTED in ", count, " ms"))

	log.Info("---------Starting integration tests---------")
	Convey("Starting integration tests", t, func() {

		Convey("Test1", func() {
			log.Info("test1")
			So(1, ShouldEqual, 1)

		})
		Convey("Test2", func() {
			log.Info("test2")
			So(1, ShouldEqual, 1)
		})
		Reset(func() {
			log.Info("reset")
		})
	})

	log.Info("Stop docktor mongoDB")
	// Stop docktor mongo container
	err = docktorMongo.Stop()
	Convey("When the container for mongoDB is stopped", t, func() {
		Convey("Then no error shoud be returned", func() {
			So(err, ShouldBeNil)
		})
	})
	log.Info("Remove docktor mongoDB")

	// Remove docktor mongo container
	err = docktorMongo.Remove(true)
	Convey("When the container for mongoDB is removed", t, func() {
		Convey("Then no error shoud be returned", func() {
			So(err, ShouldBeNil)
		})
	})
}
