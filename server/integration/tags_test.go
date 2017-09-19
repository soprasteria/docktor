package integration

import (
	"testing"

	"github.com/moby/moby/client"
	log "github.com/sirupsen/logrus"
	"github.com/soprasteria/dockerapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TagsSuite struct {
	suite.Suite
	docker  *client.Client
	mongoDB *dockerapi.Container
}

// Run all tests
func TestTagsSuite(t *testing.T) {
	suite.Run(t, new(TagsSuite))
}

// Init the docktor server and the docker daemon
func (suite *TagsSuite) SetupSuite() {
	/*
		// Start server
		log.Info("Start Server")
		StartServer()*/

	// Instantiate docker daemon
	log.Info("Instantiate docker daemon")
	docker, err := InstantiateDocker()
	suite.docker = docker
	if err != nil {
		suite.T().Fatalf("Unable to instantiate the docker daemon %v", err)
	}
}

// In the beggining of each test, instantiate and run the mongoDB instance
func (suite *TagsSuite) SetupTest() {

	if suite.docker == nil {
		suite.T().Fatalf("docker daemon is nil")
	}

	// Pull mongoDB image
	log.Info("Pull mongoDB image")
	err := PullImage(suite.docker, "mongo", "3.0")

	if err != nil {
		suite.T().Fatalf("Unable to pull the mongoDB image %v", err)
	}

	// Create mongoDB container
	log.Info("Create mongoDB container")
	res, err := CreateContainer(suite.docker, "mongo-test", "27017", "mongo", "3.0")

	if err != nil {
		suite.T().Fatalf("Unable to create the mongoDB container %v", err)
	}
	log.Info(res)

	// Start mongoDB container
	err := StartServer(suite.docker, res)
	if err != nil {
		suite.T().Fatalf("Unable to start the mongoDB container %v", err)
	}
	// Instantiate mongoDB container
	//log.Info("Instantiate mongoDB container")
	/*mongoDB, err := InstantiateContainer(suite.docker, "mongo-test", "mongo", "3.0", "27017")
	/*suite.mongoDB = mongoDB
	if err != nil {
		suite.T().Fatalf("Unable to instantiate the mongoDB container %v", err)
	}*/
	/*
		// Run docktor mongoDB
		log.Info("Run docktor mongoDB")
		err = RunContainer(suite.mongoDB)
		if err != nil {
			suite.T().Fatalf("Unable to run the mongoDB container %v", err)
		}*/
}

// In the end of each test, clean the mongoDB instance
func (suite *TagsSuite) TearDownTest() {

	/*if suite.mongoDB == nil {
		suite.T().Fatalf("mongoDB container is nil")
	}

	// Clean container mongoDB
	log.Info("Clean container mongoDB")
	err := CleanContainer(suite.mongoDB)
	if err != nil {
		suite.T().Fatalf("Unable to clean the mongoDB container %v", err)
	}*/
}

func (suite *TagsSuite) Test1() {
	log.Info("Test1")
	assert.Equal(suite.T(), 1, 1, "1 = 1")
	log.Info("Fin Test1")
}

func (suite *TagsSuite) Test2() {
	log.Info("Test2")
	assert.Equal(suite.T(), 2, 2, "2 = 2")
	log.Info("Fin Test2")
}

func (suite *TagsSuite) Test3() {
	log.Info("Test3")
	assert.Equal(suite.T(), 3, 3, "3 = 3")
	log.Info("Fin Test3")
}
