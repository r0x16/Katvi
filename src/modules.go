package src

import (
	"github.com/r0x16/Katvi/src/shared/domain"
	"github.com/r0x16/Katvi/src/shared/infraestructure/drivers"
	"github.com/r0x16/Katvi/src/shared/infraestructure/module"

	olt "github.com/r0x16/Katvi/src/olt/infraestructure/module"
)

func ProvideModules(bundle *drivers.ApplicationBundle) []domain.ApplicationModule {
	return []domain.ApplicationModule{
		&module.MainModule{Bundle: bundle},
		&olt.OltModule{Bundle: bundle},
	}
}
