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
	olt.GET("/frame/:frameId/slots", m.Bundle.ActionInjection(action.ListBoardAction))
	olt.GET("/frame/:frameId/ont", m.Bundle.ActionInjection(action.ListOntAction))
	olt.GET("/slot/:frameId/:boardId/pon", m.Bundle.ActionInjection(action.ListPonAction))
	olt.GET("/ont/:frameId/:boardId/:ponId/:ontId/ports", m.Bundle.ActionInjection(action.ListOntPorts))
}
