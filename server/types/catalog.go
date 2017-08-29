package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// CatatalogTemplate is an archetype to bootstrap a set of services
// These services are not bound together, meaning a service can later be deleted once instanciated
// When instanciated, user will choose the version of each service in the template. By default, only last versions are selected.
type CatatalogTemplate struct {
	ID          bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string          `bson:"name" json:"name"`
	Description string          `bson:"description,omitempty" json:"description,omitempty"`
	Services    []bson.ObjectId `bson:"services" json:"services"`
	Tags        []bson.ObjectId `bson:"tags" json:"tags"`
	Created     time.Time       `bson:"created" json:"created"`
	Updated     time.Time       `bson:"updated" json:"updated"`
}

// CatalogService defines a CDK service in a catalog, i.e. a service that can be deployed to a given machine
// A service contains many versions. Each version contains 1 or more containers to deploy
// Containers of a service are bound together, meaning if the service is stopped, all containers are stopped
// Jobs (checkhealth) or commands can be executed on containers
type CatalogService struct {
	ID       bson.ObjectId                    `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string                           `bson:"name" json:"name"`
	LogoPath string                           `bson:"logoPath" json:"logoPath"`
	Versions map[string]CatalogServiceVersion `bson:"versions" json:"versions"`
	Tags     []bson.ObjectId                  `bson:"tags" json:"tags"`
	Created  time.Time                        `bson:"created" json:"created"`
	Updated  time.Time                        `bson:"updated" json:"updated"`
}

// CatalogServiceVersion defines a version for a service from catalog. It contains the version number and all its metadata
// A service version is composed of many containers
type CatalogServiceVersion struct {
	// False means this version should not be used anymore, then should not be displayed as new version to deploy
	Active bool `bson:"active" json:"active"`
	// The version number of the whole service. Is usually the version of the main container
	Name string `bson:"name" json:"name"`
	// Markdown description of what has changed in this version
	Changelog string `bson:"changelog" json:"changelog"`
	// Last compatible version number for this current version
	// If empty, Docktor will not be able to upgrade automatically a previous service to this current service version
	LastCompatibleVersion string `bson:"lastCompatVersion" json:"lastCompatVersion"`
	// 1 or more container defining the service.
	// e.g. a database and a web application are two containers for a single service
	Containers []CatalogContainer `bson:"containers" json:"containers"`
	Created    time.Time          `bson:"created" json:"created"`
	Updated    time.Time          `bson:"updated" json:"updated"`
}

// CatalogContainer defines the
type CatalogContainer struct {
	// Image version like 'registryname/imagename:tag'
	Image string `bson:"image" json:"image"`
	// Available commands to execute inside container
	Commands Commands `bson:"commands" json:"commands"`
	// URls to reach the container
	URLs URLs `bson:"urls" json:"urls"`
	// Scheduled jobs for healthcheck
	Jobs Jobs `bson:"jobs" json:"jobs"`
	// Default variables for given container
	Variables Variables `bson:"variables" json:"variables"`
	// Default ports for given container
	Ports Ports `bson:"ports" json:"ports"`
	// Default volumes for given container
	Volumes Volumes `bson:"volumes" json:"volumes"`
	// Default parameters for given container
	Parameters Parameters `bson:"parameters" json:"parameters"`
}

// AddVariable adds a Variable to the Image
func (i *CatalogContainer) AddVariable(v *Variable) *CatalogContainer {
	i.Variables = append(i.Variables, *v)
	return i
}

// AddPort adds a Port to the Image
func (i *CatalogContainer) AddPort(p *Port) *CatalogContainer {
	i.Ports = append(i.Ports, *p)
	return i
}

// AddVolume adds a Volume to the Image
func (i *CatalogContainer) AddVolume(v *Volume) *CatalogContainer {
	i.Volumes = append(i.Volumes, *v)
	return i
}

// AddParameter adds a Parameter to the Image
func (i *CatalogContainer) AddParameter(p *Parameter) *CatalogContainer {
	i.Parameters = append(i.Parameters, *p)
	return i
}

// EqualsInConf checks that two catalog containers are equals in configuration
// It does not check the name for example, but will check ports, variables, parameters and volumes
func (i CatalogContainer) EqualsInConf(b CatalogContainer) bool {
	return i.Parameters.Equals(b.Parameters) &&
		i.Ports.Equals(b.Ports) &&
		i.Variables.Equals(b.Variables) &&
		i.Volumes.Equals(b.Volumes)
}

// IsIncludedInConf checks that two catalog containers are compatible in configuration
// It does not check the name for example, but will check ports, variables, parameters and volumes
func (i CatalogContainer) IsIncludedInConf(b CatalogContainer) bool {
	return i.Parameters.IsIncluded(b.Parameters) &&
		i.Ports.IsIncluded(b.Ports) &&
		i.Variables.IsIncluded(b.Variables) &&
		i.Volumes.IsIncluded(b.Volumes)
}
