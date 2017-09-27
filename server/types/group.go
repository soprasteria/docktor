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

	// Group name is used to override user-friendly name of containers
	// Meaning it has to be a valid container name
	// Here, it's a alphanum + underscore authorized string with up to 200 characters
	groupTitlePattern = `^[a-zA-Z0-9_]{1,200}$`
)

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

// Group is an entity (like a project) that gather services instances as containers
type Group struct {
	ID          bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty" validate:"required"`
	Title       string          `bson:"title" json:"title" validate:"required"`
	Description string          `bson:"description,omitempty" json:"description,omitempty"`
	FileSystems FileSystems     `bson:"filesystems" json:"filesystems" validate:"dive"`
	Members     Members         `bson:"members" json:"members" validate:"dive"`
	Tags        []bson.ObjectId `bson:"tags" json:"tags"`
	Created     time.Time       `bson:"created" json:"created"`
	Updated     time.Time       `bson:"updated" json:"updated"`
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
