package integration

import (
	"testing"

	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/soprasteria/dockerapi"
)

var (
	mongoDB *dockerapi.Container
	redis   *dockerapi.Container
)

func TestServer(t *testing.T) {

	// Instantiate docker daemon
	log.Info("Instantiate docker daemon")
	docker, err := InstantiateDocker()
	if err != nil {
		t.Fatalf("Unable to instantiate the docker daemon %v", err)
	}

	// Instantiate mongoDB container
	log.Info("Instantiate mongoDB container")
	mongoDB, err := InstantiateContainer(docker, "mongo-test", "mongo", "latest", "27017")
	if err != nil {
		t.Fatalf("Unable to instantiate the mongoDB container %v", err)
	}

	// Run docktor mongoDB
	log.Info("Run docktor mongoDB")
	err = RunContainer(mongoDB)
	if err != nil {
		t.Fatalf("Unable to run the mongoDB container %v", err)
	}

	// Instantiate redis container
	log.Info("Instantiate redis container")
	redis, err := InstantiateContainer(docker, "redis-test", "redis", "latest", "6379")
	if err != nil {
		t.Fatalf("Unable to instantiate the redis container %v", err)
	}

	// Run docktor redis
	log.Info("Run docktor redis")
	err = RunContainer(redis)
	if err != nil {
		t.Fatalf("Unable to run the redis container %v", err)
	}

	// Start server
	log.Info("Start Server")
	StartServer()

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

			// Clean containers mongoDB and redis
			log.Info("Clean containers mongoDB and redis")
			CleanContainer(mongoDB)
			log.Info("oulala")
			CleanContainer(redis)
			log.Info("	haha")

			// Instantiate mongoDB container
			log.Info("Instantiate mongoDB container")
			mongoDB, _ = InstantiateContainer(docker, "mongo-test", "mongo", "latest", "27017")

			// Run docktor mongoDB
			log.Info("Run docktor mongoDB")
			RunContainer(mongoDB)

			// Instantiate Redis container
			log.Info("Instantiate Redis container")
			redis, _ := InstantiateContainer(docker, "redis-test", "redis", "latest", "6379")

			// Run docktor redis
			log.Info("Run docktor redis")
			RunContainer(redis)
		})
	})

	// Clean containers mongoDB and redis
	log.Info("Clean containers final mongoDB and redis")
	CleanContainer(mongoDB)
	log.Info("clean mongoDB OK")
	CleanContainer(redis)
	log.Info("clean redis OK")
}
