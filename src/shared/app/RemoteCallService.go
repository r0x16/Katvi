package app

import (
	"github.com/r0x16/Katvi/src/shared/domain"
	"github.com/r0x16/Katvi/src/shared/domain/model"
)

type RemoteCallService struct {
}

func (s *RemoteCallService) Call(command *model.DeviceCommand, connector domain.DeviceConnectorProvider) (*model.CommandOutput, error) {
	if err := connector.Connect(); err != nil {
		return nil, err
	}
	defer connector.Stop()

	if err := connector.StartSession(); err != nil {
		return nil, err
	}
	defer connector.StopSession()

	exitcode := 0
	stdout, err := connector.SendCommand(command)
	if err != nil {
		if err.Error() != "timeout" {
			return nil, err
		}
		exitcode = 1
	}

	return &model.CommandOutput{
		Stdout:         stdout,
		FilteredOutput: connector.OutputFormat(stdout, command),
		ExitCode:       exitcode,
	}, nil
}
