package app

import (
	"errors"

	"github.com/r0x16/Katvi/src/olt/domain/responses"
	"github.com/r0x16/Katvi/src/shared/app"
	"github.com/r0x16/Katvi/src/shared/domain"
	"github.com/r0x16/Katvi/src/shared/domain/model"
)

type OntPortLister struct {
	Command   *model.DeviceCommand
	Connector domain.DeviceConnectorProvider
	Output    *model.CommandOutput
}

const ONTPORT_CATV_INDEX = 0
const ONTPORT_ETH_INDEX = 1
const ONTPORT_POTS_INDEX = 2

const ONTPORT_CATV_ID = 1
const ONTPORT_CATV_LINKSTATE = 2
const ONTPORT_CATV_TXPOWER = 3

const ONTPORT_ETH_ID = 1
const ONTPORT_ETH_TYPE = 2
const ONTPORT_ETH_SPEED = 3
const ONTPORT_ETH_DUPLEX = 4
const ONTPORT_ETH_LINKSTATE = 5
const ONTPORT_ETH_RINGSTATE = 6

const ONTPORT_POTS_ID = 1
const ONTPORT_POTS_PHYSICALSTATE = 2
const ONTPORT_POTS_ADMINSTATE = 3
const ONTPORT_POTS_HOOKSTATE = 4
const ONTPORT_POTS_SESSIONTYPE = 5
const ONTPORT_POTS_SERVICESTATE = 6
const ONTPORT_POTS_CALLSTATE = 7
const ONTPORT_POTS_CODEC = 8

func ListOntPorts(command *model.DeviceCommand, connector domain.DeviceConnectorProvider) *OntPortLister {
	return &OntPortLister{
		Command:   command,
		Connector: connector,
	}
}

func (lister *OntPortLister) List() (*responses.OntLocalPortsOutputCollection, error) {
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

func (lister *OntPortLister) parseOutput(output *model.CommandOutput) (*responses.OntLocalPortsOutputCollection, error) {
	ports := &responses.OntLocalPortsOutputCollection{}

	catv, err := lister.parseCatvPorts(output)
	if err != nil {
		return nil, err
	}
	ports.CatvPorts = catv

	eth, err := lister.parseEthPorts(output)
	if err != nil {
		return nil, err
	}
	ports.EthPorts = eth

	pots, err := lister.parsePotsPorts(output)
	if err != nil {
		return nil, err
	}
	ports.PotsPorts = pots

	return ports, nil
}

func (lister *OntPortLister) validateOutput(output *model.CommandOutput) error {
	if output.FilteredOutput == nil {
		return errors.New("no output found")
	}

	return nil
}

func (lister *OntPortLister) parseCatvPorts(output *model.CommandOutput) ([]*responses.CatvPort, error) {
	catv := make([]*responses.CatvPort, len(output.FilteredOutput[ONTPORT_CATV_INDEX]))

	for i, line := range output.FilteredOutput[ONTPORT_CATV_INDEX] {
		catv[i] = &responses.CatvPort{
			PortId:    line[ONTPORT_CATV_ID],
			LinkState: line[ONTPORT_CATV_LINKSTATE],
			TxPower:   line[ONTPORT_CATV_TXPOWER],
		}
	}

	return catv, nil
}

func (lister *OntPortLister) parseEthPorts(output *model.CommandOutput) ([]*responses.EthPort, error) {
	eth := make([]*responses.EthPort, len(output.FilteredOutput[ONTPORT_ETH_INDEX]))

	for i, line := range output.FilteredOutput[ONTPORT_ETH_INDEX] {
		eth[i] = &responses.EthPort{
			PortId:    line[ONTPORT_ETH_ID],
			PortType:  line[ONTPORT_ETH_TYPE],
			Speed:     line[ONTPORT_ETH_SPEED],
			Duplex:    line[ONTPORT_ETH_DUPLEX],
			LinkState: line[ONTPORT_ETH_LINKSTATE],
			RingState: line[ONTPORT_ETH_RINGSTATE],
		}
	}

	return eth, nil
}

func (lister *OntPortLister) parsePotsPorts(output *model.CommandOutput) ([]*responses.PotsPort, error) {
	pots := make([]*responses.PotsPort, len(output.FilteredOutput[ONTPORT_POTS_INDEX]))

	for i, line := range output.FilteredOutput[ONTPORT_POTS_INDEX] {
		pots[i] = &responses.PotsPort{
			PortId:        line[ONTPORT_POTS_ID],
			PhysicalState: line[ONTPORT_POTS_PHYSICALSTATE],
			AdminState:    line[ONTPORT_POTS_ADMINSTATE],
			HookState:     line[ONTPORT_POTS_HOOKSTATE],
			SessionType:   line[ONTPORT_POTS_SESSIONTYPE],
			ServiceState:  line[ONTPORT_POTS_SERVICESTATE],
			CallState:     line[ONTPORT_POTS_CALLSTATE],
			ServiceCodec:  line[ONTPORT_POTS_CODEC],
		}
	}

	return pots, nil
}
