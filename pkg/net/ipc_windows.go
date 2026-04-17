//go:build windows
// +build windows

package dnet

import (
	"errors"
	"io"
	"net"
	"time"

	"github.com/natefinch/npipe"
)

func pipePath(name string) string {
	if name == "" {
		return ""
	}
	if len(name) >= 9 && name[:9] == `\\.\\pipe\\` {
		return name
	}
	return `\\.\\pipe\\` + name
}

func MakeChannels(path string) error {
	if path == "" {
		return errors.New("path required")
	}
	return MakeIpcChannelWindows(path)
}

func MakeIpcChannelWindows(path string) error {
	pp := pipePath(path)
	l, err := npipe.Listen(pp)
	if err != nil {
		return err
	}
	return l.Close()
}

func ListenChannels(path string) (net.Listener, error) {
	if path == "" {
		return nil, errors.New("path required")
	}
	pp := pipePath(path)
	return npipe.Listen(pp)
}

func WriteInChannels(path string, data []byte) error {
	if path == "" {
		return errors.New("path required")
	}
	pp := pipePath(path)
	c, err := npipe.DialTimeout(pp, 2*time.Second)
	if err != nil {
		return err
	}
	defer c.Close()
	_, err = c.Write(data)
	return err
}

func ReadChannels(path string) ([]byte, error) {
	if path == "" {
		return nil, errors.New("path required")
	}
	pp := pipePath(path)
	c, err := npipe.DialTimeout(pp, 2*time.Second)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return io.ReadAll(c)
}

func SendInChannels(path string, data []byte) ([]byte, error) {
	if path == "" {
		return nil, errors.New("path required")
	}
	pp := pipePath(path)
	c, err := npipe.DialTimeout(pp, 2*time.Second)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	_, err = c.Write(data)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(c)
}
