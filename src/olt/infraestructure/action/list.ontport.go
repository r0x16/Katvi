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

func ListOntPorts(c echo.Context, bundle *drivers.ApplicationBundle) error {
	ont, err := parseOntPort(c)
	if err != nil {
		return err
	}

	lister := getOntPortLister(ont)

	ontPorts, err := lister.List()

	if err != nil {
		c.Logger().Error("Error listing OLT ONT ports")
		c.Logger().Error(lister.Command)
		c.Logger().Error(lister.Output)
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, errors.New("error listing ont ports"))
	}

	return c.JSON(http.StatusOK, ontPorts)
}

func getOntListPortsCommand(frameId int, boardId int, ponId int, ontId int) *model.DeviceCommand {
	return &model.DeviceCommand{
		Commands: []string{
			"enable",
			"config",
			fmt.Sprintf("interface gpon %d/%d", frameId, boardId),
			fmt.Sprintf("display ont port state %d %d catv-port all | no-more", ponId, ontId),
			fmt.Sprintf("display ont port state %d %d eth-port all | no-more", ponId, ontId),
			fmt.Sprintf("display ont port state %d %d pots-port all | no-more", ponId, ontId),
			"quit",
			"quit",
			"quit",
		},
		OutputRegex: []*regexp.Regexp{
			//                                         [OntId]       (PortId)           [PortType]      (LinkState)         (TxPower)
			regexp.MustCompile(fmt.Sprintf(`[[:blank:]]+%d[[:blank:]]+([\d]+)[[:blank:]]+CATV[[:blank:]]+([a-z]+)[[:blank:]]+(-?[\d]*\.?[\d]*)`, ontId)),
			//                                         [OntId]       (PortId)           (PortType)                 (Speed)                    (Duplex)                   (LinkState)         (RingState)
			regexp.MustCompile(fmt.Sprintf(`[[:blank:]]+%d[[:blank:]]+([\d+])[[:blank:]]+([-a-zA-Z0-9]+)[[:blank:]]+([-a-zA-Z0-9]+)[[:blank:]]+([-a-zA-Z0-9]+)[[:blank:]]+([a-z]+)[[:blank:]]+([-a-zA-Z0-9]+)`, ontId)),
			//                             (PortId)           (PhysicalState)            (AdminState)               (HookState)                (sessionType)              (ServiceState)             (CallState)                (ServiceCodec)
			regexp.MustCompile(`[[:blank:]]+([\d]+)[[:blank:]]+([-a-zA-Z0-9]+)[[:blank:]]+([-a-zA-Z0-9]+)[[:blank:]]+([-a-zA-Z0-9]+)[[:blank:]]+([-a-zA-Z0-9]+)[[:blank:]]+([-a-zA-Z0-9]+)[[:blank:]]+([-a-zA-Z0-9]+)[[:blank:]]+([-a-zA-Z0-9]+)`),
		},
		ExitRegex: regexp.MustCompile(`Check whether system data has been changed`),
		Timeout:   6 * time.Second,
	}
}

func getOntPortLister(ont map[string]int) *app.OntPortLister {
	command := getOntListPortsCommand(
		ont["frameId"],
		ont["boardId"],
		ont["ponId"],
		ont["ontId"],
	)
	connector := &devices.OLTHuaweiSshConnector{}

	return app.ListOntPorts(command, connector)
}

func parseOntPort(c echo.Context) (map[string]int, error) {
	output := make(map[string]int, 4)

	frameId, err := getBoardsFrameId(c)
	if err != nil {
		return nil, err
	}
	output["frameId"] = frameId

	boardId, err := getPonBoardId(c)
	if err != nil {
		return nil, err
	}
	output["boardId"] = boardId

	ponId, err := getOntPonId(c)
	if err != nil {
		return nil, err
	}
	output["ponId"] = ponId

	ontId, err := getOntId(c)
	if err != nil {
		return nil, err
	}
	output["ontId"] = ontId

	return output, nil
}

func getOntPonId(c echo.Context) (int, error) {
	ponId := c.Param("ponId")

	if ponId == "" {
		return 0, errors.New("invalid ponId")
	}

	return strconv.Atoi(ponId)
}

func getOntId(c echo.Context) (int, error) {
	ontId := c.Param("ontId")

	if ontId == "" {
		return 0, errors.New("invalid ontId")
	}

	return strconv.Atoi(ontId)
}
