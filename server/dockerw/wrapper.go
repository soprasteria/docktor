package dockerw

import (
	"strconv"

	"github.com/soprasteria/dockerapi"
	"github.com/soprasteria/godocktor-api/types"
)

// InitDocker : create a docker instance using daemon
func InitDocker(daemon types.Daemon) (*dockerapi.Client, error) {
	var api *dockerapi.Client
	var err error

	dockerHost := daemon.Protocol + "://" + daemon.Host + ":" + strconv.Itoa(daemon.Port)
	if daemon.Cert == "" {
		api, err = dockerapi.NewClient(daemon.Host)
	} else {
		api, err = dockerapi.NewTLSClient(dockerHost, daemon.Cert, daemon.Key, daemon.Ca)
	}
	return api, err
}