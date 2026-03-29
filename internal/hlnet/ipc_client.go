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

func (c *Client) SendCommand(cmd IPCCommand) (*IPCResponse, error) {
	if c == nil || c.conf.Path == "" {
		return nil, errors.New("invalid ipc client")
	}

	payload, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	respBytes, err := net.SendInChannels(c.conf.Path, payload)
	if err != nil {
		return nil, err
	}

	var resp IPCResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

