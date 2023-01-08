package model

import (
	"regexp"
	"time"
)

type DeviceCommand struct {
	Commands    []string
	OutputRegex []regexp.Regexp
	ExitRegex   string
	Timeout     time.Duration
}
