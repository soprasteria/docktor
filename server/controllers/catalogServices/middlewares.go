package catalogServices

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/soprasteria/docktor/server/storage"
)

// RetrieveCatalogService find catalogService using id param and put it in echo.Context
func RetrieveCatalogService(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		docktorAPI := c.Get("api").(*storage.Docktor)
		catalogServiceID := c.Param("catalogServiceID")
		if catalogServiceID == "" {
			return c.String(http.StatusBadRequest, CatalogServiceInvalidID)
		}
		catalogService, err := docktorAPI.CatalogServices().FindByID(catalogServiceID)
		if err != nil || catalogService.ID.Hex() == "" {
			log.WithField("catalogService", catalogServiceID).Error("Unable to fetch catalogService")
			return c.String(http.StatusNotFound, fmt.Sprintf(CatalogServiceNotFound, catalogServiceID))
		}
		c.Set("catalogService", catalogService)
		return next(c)
	}
}
