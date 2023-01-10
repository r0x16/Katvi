package devices

import (
	"bytes"
	"regexp"

	"github.com/r0x16/Katvi/src/shared/domain"
)

type CommandOutputWatcher struct {
	Buffer    *bytes.Buffer
	Ready     chan bool
	Connector domain.DeviceConnectorProvider
	ExitRegex *regexp.Regexp
}

func (c *CommandOutputWatcher) Write(p []byte) (n int, err error) {
	n, err = c.Buffer.Write(p)
	if err != nil {
		return n, err
	}

	// Loop through all lines in the buffer
	// If any line matches the regex, send 'true' to the 'ready' channel
	regexp := c.ExitRegex
	for _, line := range bytes.Split(c.Buffer.Bytes(), []byte("\n")) {
		if regexp.Match(line) {
			c.Ready <- true
		}
	}

	return n, err
}

func (c *CommandOutputWatcher) String() string {
	return c.Buffer.String()
}
