package module

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/r0x16/Katvi/src/shared/domain"
	"github.com/r0x16/Katvi/src/shared/infraestructure/drivers"
)

type MainModule struct {
	Bundle *drivers.ApplicationBundle
}

var _ domain.ApplicationModule = &MainModule{}

// Setups base main module routes
func (m *MainModule) Setup() {
	m.Bundle.Server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Katvi")
	})
}
