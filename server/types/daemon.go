package types

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	// HTTPProtocol is the protocol for reaching a daemon with HTTP
	HTTPProtocol DaemonProtocol = "http"
	// HTTPSProtocol is the protocol for reaching a daemon with HTTPS
	HTTPSProtocol DaemonProtocol = "https"

	// Daemon name is used to override hostname of containers
	// Meaning it has to be a valid container hostname
	// Here, it's a alphanum + underscore authorized string with up to 200 characters
	daemonNamePattern = `^[a-zA-Z0-9_]{1,200}$`
)

// Daemon defines a server where services can be deployed
type Daemon struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	// Only active daemons are monitored and can be used to deploy service
	Active bool `bson:"active" json:"active"`
	// Unique in database
	Name string `bson:"name" json:"name" validate:"required"`
	// Type of protocol to reach the Docker daemon (HTTP, HTTPS)
	Protocol DaemonProtocol `bson:"protocol" json:"protocol" validate:"required"`
	// Hostname or IP to reach the Docker daemon
	Host string `bson:"host" json:"host" validate:"required,hostname"`
	// Port to reach the Docker daemon
	Port int `bson:"port" json:"port" validate:"required,gte=0,lte=65535"`
	// Timeout in ms
	Timeout int `bson:"timeout" json:"timeout" validate:"required,gt=0"`
	// In case of using TLS with HTTPS, a ca content file is needed (see https://docs.docker.com/engine/security/https/)
	Ca string `bson:"ca" json:"ca,omitempty"`
	// In case of using TLS with HTTPS, a cert content file is needed (see https://docs.docker.com/engine/security/https/)
	Cert string `bson:"cert" json:"cert,omitempty"`
	// In case of using TLS with HTTPS, a key content file is needed (see https://docs.docker.com/engine/security/https/)
	Key string `bson:"key" json:"key,omitempty"`
	// A folder on the machine where the daemon is started.
	// It's meant to be used as default prefix for volume binding when deploying a new container.
	// Ex: MountingPoint=/data -> Volume binding=/data/GROUP1/container/a/given/path
	MountingPoint string `bson:"mountingPoint" json:"mountingPoint" validate:"required"`
	Description   string `bson:"description,omitempty" json:"description,omitempty"`
	// API endpoint of the instance of Cadvisor on the machine where the daemon is started
	// Cadvisor is used for monitoring (CPU/RAM), but also for filesystems
	CAdvisorAPI string `bson:"cadvisorApi,omitempty" json:"cadvisorApi,omitempty" validate:"omitempty,url"`
	// Localisation of the daemon
	Site bson.ObjectId `bson:"site" json:"site" validate:"required"`
	// Default container variables that will be populated at each container creation.
	// Ex: Proxy variables
	Variables Variables `bson:"variables" json:"variables"`
	// Default volume bindings that will be populated at each container creation.
	// Ex: /etc/localtime
	Volumes Volumes `bson:"volumes" json:"volumes"`
	// Tags on the daemon.
	Tags    []bson.ObjectId `bson:"tags" json:"tags"`
	Created time.Time       `bson:"created" json:"created"`
	Updated time.Time       `bson:"updated" json:"updated"`
}

// Site defines a localisation on the planet. It's meant to define where a daemon is located.
type Site struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Title     string        `bson:"title" json:"title" validate:"required"`
	Latitude  float64       `bson:"latitude" json:"latitude" validate:"required,gte=-90,lte=90"`
	Longitude float64       `bson:"longitude" json:"longitude" validate:"required,gte=-180,lte=180"`
	Created   time.Time     `bson:"created" json:"created"`
	Updated   time.Time     `bson:"updated" json:"updated"`
}

// DaemonProtocol is either HTTP or HTTPS
type DaemonProtocol string

// IsValid checks that the protocol is either HTTP or HTTPS
func (p DaemonProtocol) IsValid() bool {
	return p == HTTPProtocol || p == HTTPSProtocol
}

// AddVariable adds a Variable to the Daemon
func (d *Daemon) AddVariable(v Variable) {
	d.Variables = append(d.Variables, v)
}

// AddVolume adds a Volume to the Daemon
func (d *Daemon) AddVolume(v Volume) {
	d.Volumes = append(d.Volumes, v)
}

var daemonNameRegex = regexp.MustCompile(daemonNamePattern)

// Validate validates semantic of fields in daemon (like protocol type, pattern of the name and so on)
func (d Daemon) Validate() error {
	if !d.Protocol.IsValid() {
		return fmt.Errorf("Protocol obtained is %v, expected %v or %v", d.Protocol, HTTPProtocol, HTTPSProtocol)
	}

	if !daemonNameRegex.MatchString(d.Name) {
		return fmt.Errorf("Name %q does not match regex %q", d.Name, daemonNamePattern)
	}

	if d.Protocol == HTTPSProtocol {
		if d.Ca == "" || d.Cert == "" || d.Key == "" {
			return errors.New("Ca, Cert and Key are mandatory when using HTTPS protocol")
		}
	}

	if err := d.Variables.Validate(); err != nil {
		return err
	}

	if err := d.Volumes.Validate(); err != nil {
		return err
	}

	return nil
}

// DaemonsName returns the name of given daemons in a slice
func DaemonsName(daemons []Daemon) []string {
	res := []string{}
	for _, d := range daemons {
		res = append(res, d.Name)
	}
	return res
}
