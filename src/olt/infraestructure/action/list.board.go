package action

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/r0x16/Katvi/src/olt/app"
	"github.com/r0x16/Katvi/src/shared/domain/model"
	"github.com/r0x16/Katvi/src/shared/infraestructure/drivers"
	"github.com/r0x16/Katvi/src/shared/infraestructure/drivers/devices"
)

func ListBoardAction(c echo.Context, bundle *drivers.ApplicationBundle) error {

	frameId, err := getBoardsFrameId(c)
	if err != nil {
		c.Logger().Error("Invalid frameId")
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, errors.New("invalid frameId"))
	}

	lister := getBoardLister(frameId)

	boards, err := lister.List()
	if err != nil {
		c.Logger().Error(lister.Command)
		c.Logger().Error(lister.Output)
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, errors.New("error listing boards"))
	}

	return c.JSON(http.StatusOK, boards)
}

func getBoardLister(frameId int) *app.BoardLister {
	command := getBoardListCommand(frameId)
	connector := &devices.OLTHuaweiSshConnector{}

	return app.ListBoards(command, connector)
}

func getBoardListCommand(frameId int) *model.DeviceCommand {
	return &model.DeviceCommand{
		Commands: []string{
			"enable",
			fmt.Sprintf("display board %d | no-more", frameId),
			"quit",
		},
		OutputRegex: []*regexp.Regexp{
			// [spaces](BoardID)[spaces](BoardName)[spaces](BoardStatus)[spaces]
			regexp.MustCompile(`\r\n[[:blank:]]+([\d])[[:blank:]]+([A-Z0-9\-]+)[[:blank:]]+([a-zA-Z_]+)`),
		},
		ExitRegex: regexp.MustCompile(`Check whether system data has been changed`),
		Timeout:   4 * time.Second,
	}
}

func getBoardsFrameId(c echo.Context) (int, error) {
	frameIdParam := c.Param("frameId")

	if frameIdParam == "" {
		return 0, errors.New("frameId is required")
	}

	frameId, err := strconv.Atoi(frameIdParam)
	return frameId, err
}
