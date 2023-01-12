package app

import (
	"errors"

	"github.com/r0x16/Katvi/src/olt/domain/responses"
	"github.com/r0x16/Katvi/src/shared/app"
	"github.com/r0x16/Katvi/src/shared/domain"
	"github.com/r0x16/Katvi/src/shared/domain/model"
)

type BoardLister struct {
	Command   *model.DeviceCommand
	Connector domain.DeviceConnectorProvider
	Output    *model.CommandOutput
}

const BOARD_LIST_INDEX = 0

const SLOT_ID_INDEX = 1
const BOARD_NAME_INDEX = 2
const BOARD_STATUS_INDEX = 3

func ListBoards(command *model.DeviceCommand, connector domain.DeviceConnectorProvider) *BoardLister {
	return &BoardLister{
		Command:   command,
		Connector: connector,
	}
}

func (lister *BoardLister) List() (*responses.BoardListOutputCollection, error) {
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

func (lister *BoardLister) validateOutput(output *model.CommandOutput) error {
	if len(output.FilteredOutput) == 0 {
		return errors.New("no output returned from device")
	}

	return nil
}

func (lister *BoardLister) parseOutput(output *model.CommandOutput) (*responses.BoardListOutputCollection, error) {
	boards := &responses.BoardListOutputCollection{}
	boards.BoardListOutput = make([]*responses.BoardListOutput, len(output.FilteredOutput[BOARD_LIST_INDEX]))

	total := len(output.FilteredOutput[0])
	boards.TotalBoardCount = total

	for i, board := range output.FilteredOutput[BOARD_LIST_INDEX] {
		boards.BoardListOutput[i] = &responses.BoardListOutput{
			SlotID:      board[SLOT_ID_INDEX],
			BoardName:   board[BOARD_NAME_INDEX],
			BoardStatus: board[BOARD_STATUS_INDEX],
		}
	}

	return boards, nil
}
