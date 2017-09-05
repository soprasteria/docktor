package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	// PrimaryContainerType is the type of container considered as the main container of a service (e.g. main application)
	PrimaryContainerType ContainerType = "primary"
	// SidekickContainerType is the type of container considered as a secondary container of a service (e.g. database of main application)
	SidekickContainerType ContainerType = "sidekick"
)

// CatatalogTemplate is an archetype to bootstrap a set of services
// These services are not bound together, meaning a service can later be deleted once instanciated
// When instanciated, user will choose the version of each service in the template. By default, only last versions are selected.
type CatatalogTemplate struct {
	ID              bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
	Name            string          `bson:"name" json:"name"`
	Description     string          `bson:"description,omitempty" json:"description,omitempty"`
	CatalogServices []bson.ObjectId `bson:"catalogServices" json:"catalogServices"`
	Tags            []bson.ObjectId `bson:"tags" json:"tags"`
	Created         time.Time       `bson:"created" json:"created"`
	Updated         time.Time       `bson:"updated" json:"updated"`
}

// CatalogService defines a service in a catalog, i.e. a service that can be deployed to a given machine
// A service contains many versions. Each version contains 1 or more containers to deploy
// Containers of a service are bound together, meaning if the service is stopped, all containers are stopped
// Jobs (checkhealth) or commands can be executed on containers
type CatalogService struct {
	ID       bson.ObjectId                    `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string                           `bson:"name" json:"name"`
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

// CatalogContainer defines a container configuration inside a given catalog service
// It contains
type CatalogContainer struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	// A user-friendly name to distingish deployed containers inside a deployed service
	Name string `bson:"name" json:"name"`
	// Image version like 'registryname/imagename:tag'
	Image string `bson:"image" json:"image"`
	// Type of container (primary or sidekick), default to primary
	Type ContainerType `bson:"type" json:"type"`
	// Available commands to execute inside container
	Commands Commands `bson:"commands" json:"commands"`
	// URls to reach the container
	URLs URLs `bson:"urls" json:"urls"`
	// Scheduled job for healthcheck
	HealthCheck HealthCheck `bson:"healthCheck" json:"healthCheck"`
	// Default variables for given container
	Variables Variables `bson:"variables" json:"variables"`
	// Default ports for given container
	Ports Ports `bson:"ports" json:"ports"`
	// Default volumes for given container
	Volumes Volumes `bson:"volumes" json:"volumes"`
	// Default parameters for given container
	Parameters Parameters `bson:"parameters" json:"parameters"`
	// Default args for given container
	// Args can contain patterns to automatically fill the effective arguments like
	// For example: :
	// - Mongo connexion arguments : ['--auth', '${var1}']
	Args Args `bson:"args" json:"args"`
}

// AddVariable adds a Variable to the Image
func (i *CatalogContainer) AddVariable(v Variable) *CatalogContainer {
	i.Variables = append(i.Variables, v)
	return i
}

// AddPort adds a Port to the Image
func (i *CatalogContainer) AddPort(p Port) *CatalogContainer {
	i.Ports = append(i.Ports, p)
	return i
}

// AddVolume adds a Volume to the Image
func (i *CatalogContainer) AddVolume(v Volume) *CatalogContainer {
	i.Volumes = append(i.Volumes, v)
	return i
}

// AddParameter adds a Parameter to the Image
func (i *CatalogContainer) AddParameter(p Parameter) *CatalogContainer {
	i.Parameters = append(i.Parameters, p)
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

// ContainerType is the type of a container: 'primary' or 'sidekick'
// This type is meant to change the way containers are highlighted in display
type ContainerType string

// Command is a shell command that you can run inside the container to do an action
// This kind of command is meant to be launched by users/admin when needed
type Command struct {
	Name string `bson:"name" json:"name"`
	// Effective command to execute
	Exec string `bson:"exec" json:"exec"`
	// Arguments are parameters password to the command. By default, custom arguments are not authorized for a command
	Arguments CommandArguments `bson:"arguments,omitempty" json:"arguments,omitempty"`
	// Only members with one of these roles (or superadmin) can execute the command
	Roles   MemberRole `bson:"role" json:"role"`
	Created time.Time  `bson:"created" json:"created"`
	Updated time.Time  `bson:"updated" json:"updated"`
}

// Commands is a slice of Command
type Commands []Command

// Args are arguments passed to a given container
type Args []string

// CommandArguments are parameters passed to a command
// It's used to define arguments at runtime
type CommandArguments struct {
	// When true, arguments can be passed at runtime by user
	Authorized bool `bson:"authorized" json:"authorized"`
	// When not empty and Authorized is true, restrict the list of arguments that can used by user at runtime.
	// When empty and Authorized is true, arguments are not restricted
	RestrictedValues []string `bson:"restrictedValues,omitempty" json:"restrictedValues,omitempty"`
}

// URL for service
type URL struct {
	Label   string    `bson:"label" json:"label"`
	URL     string    `bson:"url" json:"url"`
	Created time.Time `bson:"created" json:"created"`
}

// URLs is a slice of URL
type URLs []URL

// JobType is the type of the job, defining how status are fetch (with docker exec or via http call)
type JobType string

const (
	// ExecJob is a type where status is fetched with a "Docker exec" on the container
	ExecJob JobType = "exec"
	// HTTPJob is a type where status is fetched with an HTTP call.
	HTTPJob JobType = "url"
)

// HealthCheck for service
type HealthCheck struct {
	Value       string  `bson:"value" json:"value"`       // ":internalport" if type = url, "unix command" if type= exec
	Interval    string  `bson:"interval" json:"interval"` // cron format
	Type        JobType `bson:"type" json:"type"`
	Description string  `bson:"description" json:"description"`
	Active      bool    `bson:"active" json:"active"`
}
