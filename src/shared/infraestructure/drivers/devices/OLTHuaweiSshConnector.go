package devices

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/r0x16/Katvi/src/shared/domain"
	"github.com/r0x16/Katvi/src/shared/domain/model"
	"golang.org/x/crypto/ssh"
)

type OLTHuaweiSshConnector struct {
	Host    string
	Port    string
	client  *ssh.Client
	session *ssh.Session
}

var _ domain.DeviceConnectorProvider = &OLTHuaweiSshConnector{}

// Connect implements domain.DeviceConnectorProvider
func (dc *OLTHuaweiSshConnector) Connect() error {
	sshConfig := &ssh.ClientConfig{
		User: os.Getenv("SSH_USER"),
		Auth: []ssh.AuthMethod{
			ssh.Password(os.Getenv("SSH_PASSWORD")),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         dc.Timeout(),
	}

	host := os.Getenv("SSH_HOST")
	port := os.Getenv("SSH_PORT")
	client, err := ssh.Dial("tcp", host+":"+port, sshConfig)
	if err != nil {
		return err
	}

	dc.client = client
	dc.Host = host
	dc.Port = port
	return nil
}

// SendCommand implements domain.DeviceConnectorProvider
func (dc *OLTHuaweiSshConnector) SendCommand(command *model.DeviceCommand) (string, error) {
	commands := strings.Join(command.Commands, "\n")

	ready := make(chan bool)
	buffer := &CommandOutputWatcher{
		Buffer:    &bytes.Buffer{},
		Ready:     ready,
		Connector: dc,
		ExitRegex: command.ExitRegex,
	}

	dc.session.Stdout = buffer
	defer dc.StopSession()

	err := dc.session.Start(commands)
	if err != nil {
		return "", err
	}

	for {
		select {
		case <-ready:
			return buffer.String(), nil
		case <-time.After(command.Timeout):
			return buffer.String(), errors.New("timeout")
		}
	}

}

// StartSession implements domain.DeviceConnectorProvider
func (dc *OLTHuaweiSshConnector) StartSession() error {
	session, err := dc.client.NewSession()
	if err != nil {
		return err
	}

	dc.session = session

	session.Stderr = os.Stderr

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 1000, 1000, modes); err != nil {
		return err
	}

	return nil
}

// StopSession implements domain.DeviceConnectorProvider
func (dc *OLTHuaweiSshConnector) StopSession() error {
	err := dc.session.Close()
	if err != nil {
		return err
	}

	return nil
}

func (dc *OLTHuaweiSshConnector) Stop() error {
	err := dc.client.Close()
	if err != nil {
		return err
	}

	return nil
}

// Timeout implements domain.DeviceConnectorProvider
func (*OLTHuaweiSshConnector) Timeout() time.Duration {
	return time.Second * 10
}

// OutputFormat implements domain.DeviceConnectorProvider
func (*OLTHuaweiSshConnector) OutputFormat(stdout string, command *model.DeviceCommand) [][][]string {
	// Using regex, check lines in output that match the command OutputRegex
	// and return them as a []string
	filtered := make([][][]string, len(command.OutputRegex))
	for index, regex := range command.OutputRegex {
		filtered[index] = regex.FindAllStringSubmatch(stdout, -1)
	}

	return filtered
}
