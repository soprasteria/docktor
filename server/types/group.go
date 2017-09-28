package types

import (
	"fmt"
	"regexp"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	// MemberModeratorRole is the role for user able to administrate their group
	MemberModeratorRole MemberRole = "moderator"
	// MemberUserRole is the default role for a simple user in a group
	MemberUserRole MemberRole = "member"

	// RunningStatus is the health check statu when container is up and running, without any encountered problem
	RunningStatus HealthCheckStatus = "running"
	// UnstableStatus is the health check statuswhen container is up, but needed process in it is down
	UnstableStatus HealthCheckStatus = "unstable"
	// DownStatus is the health check sttus when container is down
	DownStatus HealthCheckStatus = "down"

	// Group name is used to override user-friendly name of containers.
	// Meaning it has to be a valid container name.
	// Here, it's a alphanum + underscore authorized string with up to 200 characters
	groupTitlePattern = `^[a-zA-Z0-9_]{1,200}$`
)

// Group is an entity (like a project) that gather services instances as containers
type Group struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty" validate:"required"`
	// The name of the group. Is used to prefix automatically volume bindings and container name
	Title       string `bson:"title" json:"title" validate:"required"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	// Monitored filesystem.
	// It's meant to be used by projects for looking at space used by their tools.
	FileSystems FileSystems `bson:"filesystems" json:"filesystems" validate:"dive"`
	// Declared services on the group.
	// These services can be deployed or not
	Services Services        `bson:"services" json:"services"`
	Members  Members         `bson:"members" json:"members" validate:"dive"`
	Tags     []bson.ObjectId `bson:"tags" json:"tags"`
	Created  time.Time       `bson:"created" json:"created"`
	Updated  time.Time       `bson:"updated" json:"updated"`
}

var groupTitleRegex = regexp.MustCompile(groupTitlePattern)

// Validate validates semantic of fields in a group (like the name)
func (g Group) Validate() error {

	if !groupTitleRegex.MatchString(g.Title) {
		return fmt.Errorf("Name %q does not match regex %q", g.Title, groupTitlePattern)
	}

	if err := g.FileSystems.Validate(); err != nil {
		return err
	}

	if err := g.Members.Validate(); err != nil {
		return err
	}

	return nil
}

// NewGroup creates new group for another one.
// It helps setting default values, and cleaning duplicates
func NewGroup(g Group) Group {
	newGroup := g
	newGroup.Members = RemoveDuplicatesMember(g.Members)
	newGroup.Tags = removeDuplicatesTags(g.Tags)
	newGroup.FileSystems = removeDuplicatesFileSystems(g.FileSystems)
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
	Version          string          `bson:"version" json:"version"`
	Tags             []bson.ObjectId `bson:"tags" json:"tags"`
	Created          time.Time       `bson:"created" json:"created"`
	Updated          time.Time       `bson:"updated" json:"updated"`
}

// Services is a slice of multiple Service entities
type Services []Service

// Container is a container that belongs to a service
type Container struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	// Name of the container on the daemon
	Name string `bson:"name" json:"name"`
	// Hostname of the container on the daemon
	Hostname string `bson:"hostname" json:"hostname"`
	// Image identifies the version of the container
	Image string `bson:"image" json:"image"`
	// User-friendly name of the container, provided from catalog container
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
	// Actual args on the deployed container
	Args Args `bson:"args" json:"args"`
	// Id of the daemon where this container is deployed
	DaemonID bson.ObjectId   `bson:"daemonId,omitempty" json:"daemonId,omitempty"`
	Tags     []bson.ObjectId `bson:"tags" json:"tags"`
	Created  time.Time       `bson:"created" json:"created"`
	Updated  time.Time       `bson:"updated" json:"updated"`
}

// Containers is a slice of Container
type Containers []Container

// AddParameter adds a Parameter to the Container
func (c *Container) AddParameter(p Parameter) {
	c.Parameters = append(c.Parameters, p)
}

// AddPort adds a Port to the Container
func (c *Container) AddPort(p Port) {
	c.Ports = append(c.Ports, p)
}

// AddVariable adds a Variable to the Container
func (c *Container) AddVariable(v Variable) {
	c.Variables = append(c.Variables, v)
}

// AddVolume adds a Volume to the Container
func (c *Container) AddVolume(v Volume) {
	c.Volumes = append(c.Volumes, v)
}

// HealthCheckResult is a job result launched for the container
// It's meant to be stored in a cache like Redis
type HealthCheckResult struct {
	Message string            `bson:"message" json:"message"`
	Status  HealthCheckStatus `bson:"status" json:"status"`
}

// HealthCheckStatus is the status of a health check
type HealthCheckStatus string

// MemberRole defines the types of role available for a user as a member of a group
type MemberRole string

// IsValid checks whether the role is either moderator or member
func (role MemberRole) IsValid() bool {
	return role == MemberModeratorRole || role == MemberUserRole
}

// Member is user whois subscribed to the groupe. His role in this group defines what he is able to do.
type Member struct {
	User bson.ObjectId `bson:"user" json:"user" validate:"required"`
	Role MemberRole    `bson:"role" json:"role" validate:"required"`
}

// Validate checks if member role is valid
func (member Member) Validate() error {
	if !member.Role.IsValid() {
		return fmt.Errorf("Member role of user %v is not valid, expected 'moderator' or 'member', obtained '%v'", member.User.Hex(), member.Role)
	}
	return nil
}

// Members is a slice of multiple Member entities
type Members []Member

// RemoveDuplicatesMember from a member list
func RemoveDuplicatesMember(members Members) Members {
	result := Members{}
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
func (members Members) GetUsers() []bson.ObjectId {
	ids := []bson.ObjectId{}
	for _, m := range members {
		ids = append(ids, m.User)
	}
	return ids
}

// Validate checks if all members are valid
// Return false if at least one is not.
func (members Members) Validate() error {
	for _, m := range members {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// FileSystem is a filesystem watched by the group
type FileSystem struct {
	Daemon      bson.ObjectId `bson:"daemon" json:"daemon" validate:"required"`
	Partition   string        `bson:"partition" json:"partition" validate:"required"`
	Description string        `bson:"description,omitempty" json:"description,omitempty"`
}

// Validate checks that the filesystem is valid whenl
// Partition is not empty and does not contains the \0 and \n character
func (fs FileSystem) Validate() error {
	if !volumeNameRegex.MatchString(fs.Partition) {
		return fmt.Errorf("Partition %q does not match regex %q", fs.Partition, volumeNamePattern)
	}
	return nil
}

//FileSystems is a slice of FileSystem
type FileSystems []FileSystem

// Validate validates that all the filesystems are valid
// Returns false if at least one is not.
func (fss FileSystems) Validate() error {
	for _, fs := range fss {
		if err := fs.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// removeDuplicatesFileSystems from a filesystem list
// Duplicates are identified when using same daemon id and partition
func removeDuplicatesFileSystems(filesystems FileSystems) FileSystems {

	type key struct {
		Daemon    bson.ObjectId
		Partition string
	}

	result := FileSystems{}
	seen := map[key]bool{}
	for _, fs := range filesystems {
		fsKey := key{
			Daemon:    fs.Daemon,
			Partition: fs.Partition,
		}
		if _, ok := seen[fsKey]; !ok {
			result = append(result, fs)
			seen[fsKey] = true
		}
	}
	return result
}
