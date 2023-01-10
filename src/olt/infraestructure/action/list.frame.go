package action

import (
	"net/http"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/r0x16/Katvi/src/olt/app"
	"github.com/r0x16/Katvi/src/shared/domain/model"
	"github.com/r0x16/Katvi/src/shared/infraestructure/drivers"
	"github.com/r0x16/Katvi/src/shared/infraestructure/drivers/devices"
)

func ListFrameAction(c echo.Context, bundle *drivers.ApplicationBundle) error {

	lister := getFrameLister()

	frames, err := lister.List()
	if err != nil {
		c.Logger().Error(lister.Command)
		c.Logger().Error(lister.Output)
		c.Logger().Error(err)
		return err
	}

	return c.JSON(http.StatusOK, frames)
}

func getFrameListCommand() *model.DeviceCommand {
	return &model.DeviceCommand{
		Commands: []string{
			"enable",
			"display frame info | no-more",
			"quit",
		},
		OutputRegex: []*regexp.Regexp{
			// [spaces](FrameID)[spaces](FrameType)[spaces](FrameState)[spaces]
			regexp.MustCompile(`\s+([\d])\s+([A-Z0-9\-]+)\s+([[:alnum:]]+)\s+`),
			// [spaces]Total:[spaces](FrameCount)[spaces]
			regexp.MustCompile(`\s+Total:\s+([\d]+)\s+`),
		},
		ExitRegex: regexp.MustCompile(`Check whether system data has been changed`),
		Timeout:   4 * time.Second,
	}

}

func getFrameLister() *app.FrameLister {
	command := getFrameListCommand()
	connector := &devices.OLTHuaweiSshConnector{}

	return app.ListFrames(command, connector)
}
