package framework

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/r0x16/Katvi/src"
	"github.com/r0x16/Katvi/src/shared/domain"
	"github.com/r0x16/Katvi/src/shared/infraestructure/drivers"
)

type EchoApplicationProvider struct {
	Bundle *drivers.ApplicationBundle
}

var _ domain.ApplicationProvider = &EchoApplicationProvider{}

// Creates a new Echo server to serve http requests and response
func (app *EchoApplicationProvider) Boot() {
	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	app.Bundle.Server = server
}

// Provides the list of Echo modules to bootstrap all the routes
func (app *EchoApplicationProvider) ProvideModules() []domain.ApplicationModule {
	return src.ProvideModules(app.Bundle)
}

// Runs the HTTP server in the especified port and listens to errors
func (app *EchoApplicationProvider) Run() error {
	err := app.Bundle.Server.Start(":" + os.Getenv("GROUND_PORT"))
	// Start server
	app.Bundle.Server.Logger.Fatal(
		err,
	)

	return err
}
