package hlnet

import (
	ipcn "decaffeinated/pkg/net"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type IPCHandler func(cmd IPCCommand) (IPCResponse, error)


type Server struct {
	conf    IPCConfig
	handler IPCHandler
	ln      net.Listener
	done    chan struct{}
	mu      sync.Mutex
}

func NewServer(conf IPCConfig, handler IPCHandler) (*Server, error) {
	if conf.Path == "" {
		return nil, errors.New("ipc path required")
	}
	if handler == nil {
		return nil, errors.New("ipc handler required")
	}
	return &Server{conf: conf, handler: handler}, nil
}

func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ln != nil {
		return errors.New("server already started")
	}

	if err := ipcn.MakeChannels(s.conf.Path); err != nil {
		return err
	}
	ln, err := ipcn.ListenChannels(s.conf.Path)
	if err != nil {
		return err
	}

	s.ln = ln
	s.done = make(chan struct{})
	go s.serveLoop()
	return nil
}

func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ln == nil {
		return nil
	}
	close(s.done)
	err := s.ln.Close()
	s.ln = nil
	s.done = nil
	return err
}

func (s *Server) serveLoop() {
	for {
		select {
		case <-s.done:
			return
		default:
		}

		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.done:
				return
			default:
				log.Println("IPC accept error:", err)
				continue
			}
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(c net.Conn) {
	defer c.Close()

	body, err := io.ReadAll(c)
	if err != nil {
		log.Println("IPC read error:", err)
		return
	}

	var cmd IPCCommand
	if err := json.Unmarshal(body, &cmd); err != nil {
		writeResponse(c, IPCResponse{Status: "error", Message: "invalid json"})
		return
	}

	resp, err := s.handler(cmd)
	if err != nil {
		resp = IPCResponse{Status: "error", Message: err.Error()}
	}
	writeResponse(c, resp)
}

func writeResponse(c net.Conn, resp IPCResponse) {
	data, err := json.Marshal(resp)
	if err != nil {
		log.Println("IPC response marshal error:", err)
		return
	}
	_, _ = c.Write(data)
}
