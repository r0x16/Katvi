package model

import "time"

type DeviceCommand struct {
	Commands    []string
	OutputRegex []string
	ExitRegex   string
	Timeout     time.Duration
}
