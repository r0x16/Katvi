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

func ListPonAction(c echo.Context, bundle *drivers.ApplicationBundle) error {

	frameId, err := getBoardsFrameId(c)
	if err != nil {
		c.Logger().Error("Invalid frameId")
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, errors.New("invalid frameId"))
	}

	boardId, err := getPonBoardId(c)
	if err != nil {
		c.Logger().Error("Invalid boardId")
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, errors.New("invalid boardId"))
	}

	lister := getPonLister(frameId, boardId)
	ponList, err := lister.List()

	if err != nil {
		c.Logger().Error(lister.Command)
		c.Logger().Error(lister.Output)
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, errors.New("error listing pon"))
	}

	return c.JSON(http.StatusOK, ponList)

}

func getPonBoardId(c echo.Context) (int, error) {
	boardId := c.Param("boardId")

	if boardId == "" {
		return 0, errors.New("invalid boardId")
	}

	return strconv.Atoi(boardId)
}

func getPonLister(frameId int, boardId int) *app.PonLister {
	command := getPonListCommand(frameId, boardId)
	connector := &devices.OLTHuaweiSshConnector{}

	return app.ListPon(command, connector)
}

func getPonListCommand(frameId int, boardId int) *model.DeviceCommand {
	return &model.DeviceCommand{
		Commands: []string{
			"enable",
			fmt.Sprintf(`display board %d/%d | include (\d+\s+GPON)|(In port %d/ %d/\d+\s*, the total) | no-more`, frameId, boardId, frameId, boardId),
			"quit",
		},
		OutputRegex: []*regexp.Regexp{
			//     0     GPON        0              20             Online
			// [spaces](Port)[spaces]GPON[spaces](min-distance)[spaces](max-distance)[spaces](Status)
			regexp.MustCompile(`[[:blank:]]+([\d]+)\s+GPON\s+([\d]+)\s+([\d]+)\s+([a-zA-Z_\-]+)`),
			// In port 1/ 1/(Port), the total of ONTs are: (ONT Total), online: (ONT Online)
			regexp.MustCompile(fmt.Sprintf(`In port %d/ %d/([\d]+)\s*, the total of ONTs are:\s*([\d]+), online:\s*([\d]+)`, frameId, boardId)),
		},
		ExitRegex: regexp.MustCompile(`Check whether system data has been changed`),
		Timeout:   4 * time.Second,
	}
}
