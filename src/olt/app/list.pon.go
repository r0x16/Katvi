package app

import (
	"errors"

	"github.com/r0x16/Katvi/src/olt/domain/responses"
	"github.com/r0x16/Katvi/src/shared/app"
	"github.com/r0x16/Katvi/src/shared/domain"
	"github.com/r0x16/Katvi/src/shared/domain/model"
)

type PonLister struct {
	Command   *model.DeviceCommand
	Connector domain.DeviceConnectorProvider
	Output    *model.CommandOutput
}

const PON_LIST_INDEX = 0
const PON_ONT_LIST_INDEX = 1

const PON_ID_INDEX = 1
const PON_STATUS_INDEX = 4
const PON_MIN_DISTANCE_INDEX = 2
const PON_MAX_DISTANCE_INDEX = 3

const PON_ONT_TOTAL_INDEX = 2
const PON_ONT_ONLINE_INDEX = 3

func ListPon(command *model.DeviceCommand, connector domain.DeviceConnectorProvider) *PonLister {
	return &PonLister{
		Command:   command,
		Connector: connector,
	}
}

func (lister *PonLister) List() (*responses.PonListOutputCollection, error) {
	app := &app.RemoteCallService{}

	output, err := app.Call(lister.Command, lister.Connector)
	lister.Output = output
	if err != nil {
		return nil, err
	}

	err = lister.validateOutput(output)
	if err != nil {
		return nil, err
	}

	return lister.parseOutput(lister.Output)
}

func (lister *PonLister) validateOutput(output *model.CommandOutput) error {
	if len(output.FilteredOutput) == 0 {
		return errors.New("no output returned from device")
	}

	if len(output.FilteredOutput[PON_LIST_INDEX]) != len(output.FilteredOutput[PON_ONT_LIST_INDEX]) {
		return errors.New("PON and ONT lists are not the same length")
	}

	return nil
}

func (l *PonLister) parseOutput(output *model.CommandOutput) (*responses.PonListOutputCollection, error) {
	pon := &responses.PonListOutputCollection{}
	pon.TotalPonCount = len(output.FilteredOutput[PON_LIST_INDEX])
	pon.PonPorts = make([]*responses.PonListOutput, pon.TotalPonCount)

	for i, line := range output.FilteredOutput[PON_LIST_INDEX] {
		pon.PonPorts[i] = &responses.PonListOutput{
			PortId:      line[PON_ID_INDEX],
			PortType:    "GPON",
			PortStatus:  line[PON_STATUS_INDEX],
			MinDistance: line[PON_MIN_DISTANCE_INDEX],
			MaxDistance: line[PON_MAX_DISTANCE_INDEX],
		}
	}

	for i, line := range output.FilteredOutput[PON_ONT_LIST_INDEX] {

		if line[PON_ID_INDEX] != pon.PonPorts[i].PortId {
			return nil, errors.New("PON and ONT lists are not in sync")
		}

		pon.PonPorts[i].OntTotal = line[PON_ONT_TOTAL_INDEX]
		pon.PonPorts[i].OntOnline = line[PON_ONT_ONLINE_INDEX]
	}

	return pon, nil
}
