package daemons

const (
	// DaemonNotFound : error while retreiving daemons
	DaemonNotFound string = "Cannot find daemon %s"

	// DaemonInvalidID : ID is not set in url or is invalid
	DaemonInvalidID string = "Invalid daemon ID"

	// DaemonInvalidInfo : information about docker daemon is invalid
	DaemonInvalidInfo string = "Invalid daemon info"

	// daemonInfoCacheKeyFormat is the key format to store Docker stats of a daemon in a cache.
	// %s in format is meant to be the daemon id
	daemonInfoCacheKeyFormat string = "daemon:%s:info"
)
