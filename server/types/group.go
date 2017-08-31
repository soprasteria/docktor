package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	// MemberModeratorRole is the role for user able to administrate their group
	MemberModeratorRole MemberRole = "moderator"
	// MemberUserRole is the default role for a simple user in a group
	MemberUserRole MemberRole = "member"
	// RunningStatus is the health check status when container is up and running, without any encountered problem
	RunningStatus HealthCheckStatus = "running"
	// UnstableStatus is the health check status when container is up, but needed process in it is down
	UnstableStatus HealthCheckStatus = "unstable"
	// DownStatus is the health check status when container is down
	DownStatus HealthCheckStatus = "down"
)

// Group is a entity (like a project) that gather services instances as containers
type Group struct {
	ID          bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
	Created     time.Time       `bson:"created" json:"created"`
	Title       string          `bson:"title" json:"title"`
	Description string          `bson:"description" json:"description"`
	FileSystems FileSystems     `bson:"filesystems" json:"filesystems"`
	Services    Services        `bson:"services" json:"services"`
	Members     Members         `bson:"members" json:"members"`
	Tags        []bson.ObjectId `bson:"tags" json:"tags"`
}

// NewGroup creates new group for another one.
// It helps setting default values, and cleaning duplicates
func NewGroup(g Group) Group {
	newGroup := g
	newGroup.Members = RemoveDuplicatesMember(g.Members)
	newGroup.Tags = removeDuplicatesTags(g.Tags)
	return newGroup
}

// AddFileSystem adds a FileSystem to the Group
func (g *Group) AddFileSystem(f FileSystem) {
	g.FileSystems = append(g.FileSystems, f)
}

// AddService adds a Service to the Group
func (g *Group) AddService(s Service) {
	g.Services = append(g.Services, s)
}

// Service is a deployed Service in a group.
// It's a virtual entity composed of multiple containers
type Service struct {
	ID               bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
	Name             string          `bson:"name" json:"name"`
	Description      string          `bson:"description,omitempty" json:"description,omitempty"`
	CatalogServiceID bson.ObjectId   `bson:"catalogServiceId" json:"catalogServiceId"`
	version          string          `bson:"version" json:"version"`
	Tags             []bson.ObjectId `bson:"tags" json:"tags"`
}

// Services is a slice of multiple Service entities
type Services []Service

// Container is a container associated to the group
type Container struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	// Name of the container on the daemon
	Name string `bson:"name" json:"name"`
	// Hostname of the container on the daemon
	Hostname string `bson:"hostname" json:"hostname"`
	// Image identifies the version of the container
	Image string `bson:"image" json:"image"`
	// Name of the container type, from catalog
	CatalogContainerName string `bson:"catalogContainerName" json:"catalogContainerName"`
	// Id of the container type, from catalog
	CatalogContainerID bson.ObjectId `bson:"catalogContainerId" json:"catalogContainerId"`
	// Full id of the container on the daemon
	ContainerID string `bson:"containerId" json:"containerId"`
	// Actual parameters used on the deployed container
	Parameters Parameters `bson:"parameters" json:"parameters"`
	// Actual bound ports on the deployed container
	Ports Ports `bson:"ports" json:"ports"`
	// Actual bound variables on the deployed container
	Variables Variables `bson:"variables" json:"variables"`
	// Actual volumes mapped on the deployed container
	Volumes Volumes `bson:"volumes" json:"volumes"`
	// Results of healthcheck jobs executed on this kind of containers
	HealthCheck HealthCheckResults `bson:"healthCheck" json:"healthCheck"`
	// Id of the daemon where this container is deployed
	DaemonID bson.ObjectId   `bson:"daemonId,omitempty" json:"daemonId,omitempty"`
	Tags     []bson.ObjectId `bson:"tags" json:"tags"`
}

// Containers is a slice of Container
type Containers []Container

// AddParameter adds a ParameterContainer to the Container
func (c *Container) AddParameter(p Parameter) {
	c.Parameters = append(c.Parameters, p)
}

// AddPort adds a PortContainer to the Container
func (c *Container) AddPort(p Port) {
	c.Ports = append(c.Ports, p)
}

// AddVariable adds a VariableContainer to the Container
func (c *Container) AddVariable(v Variable) {
	c.Variables = append(c.Variables, v)
}

// AddVolume adds a VolumeContainer to the Container
func (c *Container) AddVolume(v Volume) {
	c.Volumes = append(c.Volumes, v)
}

// HealthCheckResults is a job lunched for the container
type HealthCheckResults struct {
	Message string            `bson:"message" json:"message"`
	Status  HealthCheckStatus `bson:"status" json:"status"`
}

// HealthCheckStatus is the status of a health check
//
type HealthCheckStatus string

// MemberRole defines the types of role available for a user as a member of a group
type MemberRole string

// Member is user whois subscribed to the groupe. His role in this group defines what he is able to do.
type Member struct {
	User bson.ObjectId `bson:"user" json:"user"`
	Role MemberRole    `bson:"role" json:"role"`
}

// Members is a slice of multiple Member entities
type Members []Member

// RemoveDuplicatesMember from a member list
func RemoveDuplicatesMember(members Members) Members {
	var result Members
	seen := map[bson.ObjectId]bool{}
	for _, member := range members {
		if _, ok := seen[member.User]; !ok {
			result = append(result, member)
			seen[member.User] = true
		}
	}
	return result
}

//GetUsers gets ids of members
func (members *Members) GetUsers() []bson.ObjectId {
	ids := []bson.ObjectId{}
	for _, m := range *members {
		ids = append(ids, m.User)
	}
	return ids
}

// FileSystem is a filesystem watched by the group
type FileSystem struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Daemon      bson.ObjectId `bson:"daemon" json:"daemon"`
	Partition   string        `bson:"partition,omitempty" json:"partition,omitempty"`
	Description string        `bson:"description" json:"description"`
}

//FileSystems is a slice of FileSystem
type FileSystems []FileSystem

// ContainerWithGroup is a entity which contains a container, linked to a group
type ContainerWithGroup struct {
	Group     Group
	Container Container
}

// ContainerWithGroupID is an entity which contains a container, linked to a group ID
type ContainerWithGroupID struct {
	Container Container     `bson:"container" json:"container"`
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
}
