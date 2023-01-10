package app

import (
	"errors"
	"strconv"

	"github.com/r0x16/Katvi/src/olt/domain/responses"
	"github.com/r0x16/Katvi/src/shared/app"
	"github.com/r0x16/Katvi/src/shared/domain"
	"github.com/r0x16/Katvi/src/shared/domain/model"
)

type FrameLister struct {
	Command   *model.DeviceCommand
	Connector domain.DeviceConnectorProvider
	Output    *model.CommandOutput
}

const FRAME_LIST_INDEX = 0
const FRAME_COUNT_INDEX = 1

const FRAME_ID_INDEX = 1
const FRAME_TYPE_INDEX = 2
const FRAME_STATE_INDEX = 3

func ListFrames(command *model.DeviceCommand, connector domain.DeviceConnectorProvider) *FrameLister {
	return &FrameLister{
		Command:   command,
		Connector: connector,
	}
}

func (lister *FrameLister) List() (*responses.FrameListOutputCollection, error) {
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

func (lister *FrameLister) parseOutput(output *model.CommandOutput) (*responses.FrameListOutputCollection, error) {
	frames := &responses.FrameListOutputCollection{}
	frames.FrameListOutput = make([]responses.FrameListOutput, len(output.FilteredOutput[0]))
	total, err := strconv.Atoi(output.FilteredOutput[FRAME_COUNT_INDEX][0][1])

	if err != nil {
		return nil, err
	}

	frames.TotalFrameCount = total

	for i, frame := range output.FilteredOutput[0] {
		frames.FrameListOutput[i].FrameID = frame[FRAME_ID_INDEX]
		frames.FrameListOutput[i].FrameType = frame[FRAME_TYPE_INDEX]
		frames.FrameListOutput[i].FrameState = frame[FRAME_STATE_INDEX]
	}

	return frames, nil
}

func (lister *FrameLister) validateOutput(output *model.CommandOutput) error {
	if output.FilteredOutput == nil {
		return errors.New("no output found")
	}

	if len(output.FilteredOutput[FRAME_LIST_INDEX]) == 0 {
		return errors.New("no frames found")
	}

	if len(output.FilteredOutput[FRAME_COUNT_INDEX]) != 1 {
		return errors.New("no frame count found")
	}

	total, err := strconv.Atoi(output.FilteredOutput[FRAME_COUNT_INDEX][0][1])
	if err != nil {
		return errors.New("error parsing frame count")
	}

	if len(output.FilteredOutput[0]) != total {
		return errors.New("frame count mismatch")
	}

	return nil
}
