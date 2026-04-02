package hlnet

import (
	"decaffeinated/pkg/net"
	"encoding/json"
	"errors"
)

// IPCConfig carries IPC endpoint settings.
type IPCConfig struct {
	Path string
}

// IPCCommand is synced with WatchDog IPC command format.
type IPCCommand struct {
	Action           string `json:"action"`
	AppName          string `json:"app_name"`
	TimeLimitSeconds int64  `json:"time_limit_seconds,omitempty"`
	IsBlocked        bool   `json:"is_blocked,omitempty"`
}

type CommandBundle struct {
	Commands []IPCCommand
}
// IPCResponse is returned by the WatchDog IPC endpoint.
type IPCResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Client holds IPC destination information.
type Client struct {
	conf IPCConfig
}

func NewClient(conf IPCConfig) (*Client, error) {
	if conf.Path == "" {
		return nil, errors.New("ipc path required")
	}
	return &Client{conf: conf}, nil
}

func (c *Client) SendCommand(cmd IPCCommand) error {
	if c == nil || c.conf.Path == "" {
		return errors.New("invalid ipc client")
	}

	payload, err := json.Marshal(cmd)
	if err != nil {
		return nil
	}

	err = net.WriteInChannels(c.conf.Path, payload)
	if err != nil {
		return nil
	}
	return err
}

