package app

import (
	"fmt"
	"strconv"

	"github.com/r0x16/Katvi/src/olt/domain/responses"
	"github.com/r0x16/Katvi/src/shared/app"
	"github.com/r0x16/Katvi/src/shared/domain"
	"github.com/r0x16/Katvi/src/shared/domain/model"
)

type OntLister struct {
	Command   *model.DeviceCommand
	Connector domain.DeviceConnectorProvider
	Output    *model.CommandOutput
}

const ONT_STATE_LIST_INDEX = 0
const ONT_UPTIME_LIST_INDEX = 1
const ONT_DETAILS_LIST_INDEX = 2
const ONT_TOTALS_INDEX = 3
const ONT_SERVICES_INDEX = 4

const ONT_ID_INDEX = 1
const ONT_SERIALNUMBER_INDEX = 2
const ONT_CONTROLFLAG_INDEX = 3
const ONT_RUNSTATE_INDEX = 4
const ONT_CONFIGSTATE_INDEX = 5
const ONT_MATCHSTATE_INDEX = 6
const ONT_PROTECTSIDE_INDEX = 7

const ONT_UPTIME_INDEX = 2
const ONT_DOWNTIME_INDEX = 3
const ONT_DOWNCASE_INDEX = 4

const ONT_TYPE_INDEX = 2
const ONT_DISTANCE_INDEX = 3
const ONT_RXPOWER_INDEX = 4
const ONT_TXPOWER_INDEX = 5
const ONT_DESCRIPTION_INDEX = 6

const ONT_TOTAL_FRAME_INDEX = 1
const ONT_TOTAL_BOARD_INDEX = 2
const ONT_TOTAL_PON_INDEX = 3
const ONT_TOTAL_INDEX = 4
const ONT_ONLINE_INDEX = 5

const ONT_SERVICE_ID_INDEX = 1
const ONT_SERVICE_VLAN_INDEX = 2
const ONT_SERVICE_VLANATTR_INDEX = 3
const ONT_SERVICE_PORTTYPE_INDEX = 4
const ONT_SERVICE_FRAME_INDEX = 5
const ONT_SERVICE_BOARD_INDEX = 6
const ONT_SERVICE_PON_INDEX = 7
const ONT_SERVICE_ONT_INDEX = 8
const ONT_SERVICE_GEMPORT_INDEX = 9
const ONT_SERVICE_FLOWTYPE_INDEX = 10
const ONT_SERVICE_RX_INDEX = 11
const ONT_SERVICE_TX_INDEX = 12
const ONT_SERVICE_STATE_INDEX = 13

func ListOnt(command *model.DeviceCommand, connector domain.DeviceConnectorProvider) *OntLister {
	return &OntLister{
		Command:   command,
		Connector: connector,
	}
}

func (lister *OntLister) List() (*responses.OntListOutputCollection, error) {
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

func (lister *OntLister) validateOutput(output *model.CommandOutput) error {
	if len(output.FilteredOutput) != 5 {
		return fmt.Errorf("no output returned from device")
	}

	if len(output.FilteredOutput[ONT_TOTALS_INDEX]) == 0 {
		return fmt.Errorf("total ONT count not returned from device")
	}

	if len(output.FilteredOutput[ONT_STATE_LIST_INDEX]) != len(output.FilteredOutput[ONT_UPTIME_LIST_INDEX]) {
		return fmt.Errorf("ONT state and uptime lists are not the same length")
	}

	if len(output.FilteredOutput[ONT_STATE_LIST_INDEX]) != len(output.FilteredOutput[ONT_DETAILS_LIST_INDEX]) {
		return fmt.Errorf("ONT state and details lists are not the same length")
	}

	return nil
}

func (lister *OntLister) parseOutput(output *model.CommandOutput) (*responses.OntListOutputCollection, error) {
	ontOutput := &responses.OntListOutputCollection{}
	ontOutput.OntLists = make(map[string]*responses.OntByPonList)
	ontOutput.TotalOntCount = 0
	ontOutput.OnlineOntCount = 0
	ontOutput.OfflineOntCount = 0

	offset := 0

	for _, ponTotals := range output.FilteredOutput[ONT_TOTALS_INDEX] {
		frame := ponTotals[ONT_TOTAL_FRAME_INDEX]
		board := ponTotals[ONT_TOTAL_BOARD_INDEX]
		pon := ponTotals[ONT_TOTAL_PON_INDEX]

		key := fmt.Sprintf("%s/%s/%s", frame, board, pon)

		onts := &responses.OntByPonList{}
		err := lister.setTotals(onts, ponTotals)
		if err != nil {
			return nil, err
		}

		lister.setPonOnts(offset, onts, output)

		ontOutput.OntLists[key] = onts
		ontOutput.TotalOntCount += onts.TotalOntCount
		ontOutput.OnlineOntCount += onts.OnlineOntCount
		ontOutput.OfflineOntCount += onts.OfflineOntCount

		offset += onts.TotalOntCount

	}

	lister.setOntServices(ontOutput, output)

	return ontOutput, nil
}

func (lister *OntLister) setTotals(onts *responses.OntByPonList, totals []string) error {
	total, err := strconv.Atoi(totals[ONT_TOTAL_INDEX])
	if err != nil {
		return err
	}

	online, err := strconv.Atoi(totals[ONT_ONLINE_INDEX])
	if err != nil {
		return err
	}

	onts.TotalOntCount = total
	onts.OnlineOntCount = online
	onts.OfflineOntCount = total - online

	return nil

}

func (lister *OntLister) setPonOnts(offset int, onts *responses.OntByPonList, output *model.CommandOutput) (*responses.OntByPonList, error) {
	onts.Onts = make(map[string]*responses.OntListOutput)

	for i := offset; i < (offset + onts.TotalOntCount); i++ {

		stateId := output.FilteredOutput[ONT_STATE_LIST_INDEX][i][ONT_ID_INDEX]
		uptimeId := output.FilteredOutput[ONT_UPTIME_LIST_INDEX][i][ONT_ID_INDEX]
		detailsId := output.FilteredOutput[ONT_DETAILS_LIST_INDEX][i][ONT_ID_INDEX]

		if stateId != uptimeId || stateId != detailsId {
			return nil, fmt.Errorf("ONT state, uptime and details lists are not in sync")
		}

		onuId := output.FilteredOutput[ONT_STATE_LIST_INDEX][i][ONT_ID_INDEX]

		onts.Onts[onuId] = &responses.OntListOutput{
			Id:           onuId,
			SerialNumber: output.FilteredOutput[ONT_STATE_LIST_INDEX][i][ONT_SERIALNUMBER_INDEX],
			ControlFlag:  output.FilteredOutput[ONT_STATE_LIST_INDEX][i][ONT_CONTROLFLAG_INDEX],
			RunState:     output.FilteredOutput[ONT_STATE_LIST_INDEX][i][ONT_RUNSTATE_INDEX],
			ConfigState:  output.FilteredOutput[ONT_STATE_LIST_INDEX][i][ONT_CONFIGSTATE_INDEX],
			MatchState:   output.FilteredOutput[ONT_STATE_LIST_INDEX][i][ONT_MATCHSTATE_INDEX],
			ProtectSide:  output.FilteredOutput[ONT_STATE_LIST_INDEX][i][ONT_PROTECTSIDE_INDEX],
			LastUpTime:   output.FilteredOutput[ONT_UPTIME_LIST_INDEX][i][ONT_UPTIME_INDEX],
			LastDownTime: output.FilteredOutput[ONT_UPTIME_LIST_INDEX][i][ONT_DOWNTIME_INDEX],
			LastDownCase: output.FilteredOutput[ONT_UPTIME_LIST_INDEX][i][ONT_DOWNCASE_INDEX],
			Type:         output.FilteredOutput[ONT_DETAILS_LIST_INDEX][i][ONT_TYPE_INDEX],
			Distance:     output.FilteredOutput[ONT_DETAILS_LIST_INDEX][i][ONT_DISTANCE_INDEX],
			RxPower:      output.FilteredOutput[ONT_DETAILS_LIST_INDEX][i][ONT_RXPOWER_INDEX],
			TxPower:      output.FilteredOutput[ONT_DETAILS_LIST_INDEX][i][ONT_TXPOWER_INDEX],
			Description:  output.FilteredOutput[ONT_DETAILS_LIST_INDEX][i][ONT_DESCRIPTION_INDEX],
			Services:     make(map[string]*responses.OntService),
		}
	}

	return onts, nil
}

func (lister *OntLister) setOntServices(ontOutput *responses.OntListOutputCollection, output *model.CommandOutput) error {
	services := output.FilteredOutput[ONT_SERVICES_INDEX]

	for _, service := range services {

		ontId := service[ONT_SERVICE_ONT_INDEX]
		pon := service[ONT_SERVICE_PON_INDEX]
		board := service[ONT_SERVICE_BOARD_INDEX]
		frame := service[ONT_SERVICE_FRAME_INDEX]

		key := fmt.Sprintf("%s/%s/%s", frame, board, pon)

		onts := ontOutput.OntLists[key]

		if onts == nil {
			return fmt.Errorf("ONT list for PON %s not found", key)
		}

		ont := onts.Onts[ontId]

		if ont == nil {
			return fmt.Errorf("ONT %s not found on PON %s", ontId, key)
		}

		ont.Services[service[ONT_SERVICE_ID_INDEX]] = &responses.OntService{
			Index:    service[ONT_SERVICE_ID_INDEX],
			VlanId:   service[ONT_SERVICE_VLAN_INDEX],
			VlanAttr: service[ONT_SERVICE_VLANATTR_INDEX],
			PortType: service[ONT_SERVICE_PORTTYPE_INDEX],
			GemPort:  service[ONT_SERVICE_GEMPORT_INDEX],
			FlowType: service[ONT_SERVICE_FLOWTYPE_INDEX],
			Rx:       service[ONT_SERVICE_RX_INDEX],
			Tx:       service[ONT_SERVICE_TX_INDEX],
			State:    service[ONT_SERVICE_STATE_INDEX],
		}
	}

	return nil
}
