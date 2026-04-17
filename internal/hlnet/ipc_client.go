package hlnet

import (
	"net"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

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


