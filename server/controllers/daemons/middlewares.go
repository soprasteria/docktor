package daemons

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/soprasteria/docktor/server/storage"
	"github.com/spf13/viper"
)

// RetrieveDaemon find daemon using id param and put it in echo.Context
func RetrieveDaemon(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		docktorAPI := c.Get("api").(*storage.Docktor)
		daemonID := c.Param("daemonID")
		if daemonID == "" {
			return c.String(http.StatusBadRequest, DaemonInvalidID)
		}
		daemon, err := docktorAPI.Daemons().FindByID(daemonID)
		if err != nil {
			log.WithError(err).Errorf("Unable to find given daemon %v", daemon.ID)
			return c.String(http.StatusNotFound, fmt.Sprintf(DaemonNotFound, daemonID))
		}
		daemon, err = DecryptDaemon(daemon, viper.GetString("auth.encrypt-secret"))
		if err != nil {
			log.WithError(err).Errorf("Unable to decrypt daemon %v", daemon.ID)
			return c.String(http.StatusInternalServerError, "Unable to get daemon because of technical error. Retry later.")
		}
		c.Set("daemon", daemon)
		return next(c)
	}
}
