package domain

import (
	"time"

	"github.com/r0x16/Katvi/src/shared/domain/model"
)

type DeviceConnectorProvider interface {
	Connect() error
	StartSession() error
	StopSession() error
	Timeout() time.Duration
	SendCommand(model.DeviceCommand) (string, error)
}
