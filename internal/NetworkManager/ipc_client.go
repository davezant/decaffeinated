package hlnet

import (
	"net"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	DefaultLinuxSockPath = `./sock`
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

type CommandBundle struct {
	Commands []IPCCommand
}

type Client struct {
	conf IPCConfig
}

func NewClient(conf IPCConfig) (*Client, error) {
	if conf.Path == "" {
		return nil, errors.New("ipc path required")
	}
	return &Client{conf: conf}, nil
}

func (c *Client) SendCommand(cmd IPCCommand) (IPCResponse, error) {

    conn, err := net.DialTimeout("unix", c.conf.Path, 5*time.Second)
    if err != nil {
        return IPCResponse{}, fmt.Errorf("failed to connect to watchdog: %w", err)
    }
    defer conn.Close()

    if err := json.NewEncoder(conn).Encode(cmd); err != nil {
        return IPCResponse{}, fmt.Errorf("failed to encode command: %w", err)
    }

    var resp IPCResponse
    if err := json.NewDecoder(conn).Decode(&resp); err != nil {
        return IPCResponse{}, fmt.Errorf("failed to decode response: %w", err)
    }

    return resp, nil
}


