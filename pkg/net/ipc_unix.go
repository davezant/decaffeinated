//go:build !windows
// +build !windows

package dnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func MakeChannels(path string) error {
	if path == "" {
		return errors.New("path required")
	}
	return MakeIpcChannelLinux(path)
}

func MakeIpcChannelLinux(path string) error {
	if path == "" {
		return errors.New("path required")
	}
	_ = os.Remove(path)
	l, err := net.Listen("unix", path)
	if err != nil {
		return err
	}
	return l.Close()
}

func ListenChannels(path string) (net.Listener, error) {
	if path == "" {
		return nil, errors.New("path required")
	}
	_ = os.Remove(path)

	return net.Listen("unix", path)
}

func WriteInChannels(path string, data []byte) error {
	if path == "" {
		return errors.New("path required")
	}
	c, err := net.DialTimeout("unix", path, 2*time.Second)
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
	c, err := net.DialTimeout("unix", path, 20* time.Second)
	fmt.Println()
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
	WriteInChannels(path ,data)
	time.Sleep(2 * time.Second)
	c, err := ReadChannels(path)
	if err != nil {
		return nil, err
	}
	return c, nil
}
