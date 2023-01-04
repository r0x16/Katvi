package infraestructure

import (
	"github.com/r0x16/Katvi/src/shared/app"
	"github.com/r0x16/Katvi/src/shared/infraestructure/drivers"
	"github.com/r0x16/Katvi/src/shared/infraestructure/drivers/db"
	"github.com/r0x16/Katvi/src/shared/infraestructure/drivers/framework"
)

type Main struct {
	mainService *app.MainService
}

func (m *Main) RunServices() {
	m.mainService = &app.MainService{}

	// Creates a new application bundle
	dbProvider := &db.GormMysqlDatabaseProvider{}
	app := &framework.EchoApplicationProvider{
		Bundle: &drivers.ApplicationBundle{
			Database: dbProvider,
		},
	}

	// Prepares the application bundle
	err := m.mainService.Run(
		app,
		dbProvider,
	)

	// Captures application error, must be logged.
	if err != nil {
		panic(err)
	}
}
