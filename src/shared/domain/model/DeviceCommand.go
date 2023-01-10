package model

import (
	"regexp"
	"time"
)

type DeviceCommand struct {
	Commands    []string
	OutputRegex []*regexp.Regexp
	ExitRegex   *regexp.Regexp
	Timeout     time.Duration
}
