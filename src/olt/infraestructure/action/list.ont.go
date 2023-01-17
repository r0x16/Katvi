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

func ListOntAction(c echo.Context, bundle *drivers.ApplicationBundle) error {
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

	ponId, err := getOntPonId(c)
	if err != nil {
		c.Logger().Error("Invalid ponId")
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, errors.New("invalid ponId"))
	}

	lister := getOntLister(frameId, boardId, ponId)
	onts, err := lister.List()

	if err != nil {
		c.Logger().Error("Error listing OLT ONTs")
		c.Logger().Error(lister.Command)
		c.Logger().Error(lister.Output)
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, errors.New("error listing ont"))
	}

	return c.JSON(http.StatusOK, onts)
}

func getOntPonId(c echo.Context) (int, error) {
	ponId := c.Param("ponId")

	if ponId == "" {
		return 0, errors.New("invalid ponId")
	}

	return strconv.Atoi(ponId)
}

func getOntLister(frameId int, boardId int, ponId int) *app.OntLister {
	command := getOntListCommand(frameId, boardId, ponId)
	connector := &devices.OLTHuaweiSshConnector{}

	return app.ListOnt(command, connector)
}

func getOntListCommand(frameId int, boardId int, ponId int) *model.DeviceCommand {
	return &model.DeviceCommand{
		Commands: []string{
			"enable",
			"config",
			fmt.Sprintf("interface gpon %d/%d", frameId, boardId),
			fmt.Sprintf("display ont info %d all | no-more", ponId),
			fmt.Sprintf("display ont info summary %d | no-more", ponId),
			"quit",
			"quit",
			"quit",
		},
		OutputRegex: []*regexp.Regexp{
			//  [frame]/ [board]/[pon] (ONT ID) (S/N) (Control Flag) (Run state) (Config State) (Match State) (Protect Side)
			regexp.MustCompile(`[[:blank:]]+[\d]+/[[:blank:]]*[\d]+/[\d]+[[:blank:]]+([\d]+)[[:blank:]]+([A-Z0-9]+)[[:blank:]]+([a-z]+)[[:blank:]]+([a-z]+)[[:blank:]]+([a-z]+)[[:blank:]]+([a-z]+)[[:blank:]]+([a-z]+)`),
			//  (ONT ID) [Run state] (Last UpTime) (Last DownTime) (Last DownCase)
			regexp.MustCompile(`[[:blank:]]+([\d]+)[[:blank:]]+[a-z]+[[:blank:]]+(-|[\d]{4}-[\d]{2}-[\d]{2} [\d]{2}:[\d]{2}:[\d]{2})[[:blank:]]+(-|[\d]{4}-[\d]{2}-[\d]{2} [\d]{2}:[\d]{2}:[\d]{2})[[:blank:]]+([[:graph:]]+)`),
			//	(ONT ID) [S/N] (Type) (Distance) (RxPower)/(TxPower) (Description)
			regexp.MustCompile(`[[:blank:]]+([\d]+)[[:blank:]]+[A-Z0-9]+[[:blank:]]+([A-Za-z0-9\-]+)[[:blank:]]+([0-9\-]+)[[:blank:]]+(-?[\d]*\.?[\d]*)\/(-?[\d]*\.?[\d]*)[[:blank:]]+([[:graph:]]+)`),
			// In port [frame]/ [board]/[pon], the total of ONTs are: (ONT Total), online: (ONT Online)
			regexp.MustCompile(fmt.Sprintf(`In port %d/ %d/%d\s*, the total of ONTs are:\s*([\d]+), online:\s*([\d]+)`, frameId, boardId, ponId)),
		},
		ExitRegex: regexp.MustCompile(`Check whether system data has been changed`),
		Timeout:   4 * time.Second,
	}
}
