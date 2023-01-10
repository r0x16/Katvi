package module

import (
	"github.com/r0x16/Katvi/src/olt/infraestructure/action"
	"github.com/r0x16/Katvi/src/shared/domain"
	"github.com/r0x16/Katvi/src/shared/infraestructure/drivers"
)

type OltModule struct {
	Bundle *drivers.ApplicationBundle
}

var _ domain.ApplicationModule = &OltModule{}

// Setup implements domain.ApplicationModule
func (m *OltModule) Setup() {
	olt := m.Bundle.Server.Group("/olt")

	olt.GET("/frame", m.Bundle.ActionInjection(action.ListFrameAction))
	olt.GET("/frame/:frameId/boards", m.Bundle.ActionInjection(action.ListBoardAction))
}
