package hlnet

import (
	"decaffeinated/pkg/net"
	"encoding/json"
	"errors"
)

const (
	DefaultLinuxSockPath = `/tmp/decaffeinated-socket`
	DefaultWindowsPipePath = `\\.\pipe\decaffeinated`
)

type IPCConfig struct {
	Path string
}

type IPCCommand struct {
	Action           string `json:"action"`
	AppNames          []string `json:"app_names"`
	Category		 string `json:"category"`
	TimeLimitSeconds int64  `json:"time_limit_seconds,omitempty"`
	IsBlocked        bool   `json:"is_blocked,omitempty"`
	CustomTimestamps []CustomTimestamp `json:"rule,omitempty"`
}
 
type CustomTimestamp struct {
	Timestamp float32 `json:"timestamp"`
	Callback  string  `json:"callback"`
}

type CommandBundle struct {
	Commands []IPCCommand
}

// Heartbeat for external User control.
type Heartbeat struct {
	UserUUID string `json:"UUID"`
	TimeLimitSeconds int64 `json:"time_limit_seconds,omitempty"`
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
	// TODO WAIT UNTIL HEAR
	return err
}

func (c *Client) WaitUntilResponse(){

}
