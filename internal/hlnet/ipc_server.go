package hlnet

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

// IPCHandler define a assinatura da função que processará os comandos
type IPCHandler func(cmd IPCCommand) (IPCResponse, error)



type Server struct {
    conf    IPCConfig
    handler IPCHandler
    ln      net.Listener
    done    chan struct{}
    mu      sync.Mutex
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
        return
    }
    c.Write(data)
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

// O método deve ser Start (Maiúsculo)
func (s *Server) Start() error {
	if s.ln != nil {
		return errors.New("server already started")
	}

	// Substitua pela sua lógica de socket real
	ln, err := net.Listen("unix", s.conf.Path)
	if err != nil {
		return err
	}

	s.ln = ln
	s.done = make(chan struct{})
	
	go s.serveLoop() // Certifique-se que serveLoop existe no mesmo arquivo
	return nil
}

// O método deve ser Stop (Maiúsculo)
func (s *Server) Stop() error {
	if s.ln == nil {
		return nil
	}
	
	close(s.done)
	err := s.ln.Close()
	s.ln = nil
	return err
}

func (s *Server) serveLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.done:
				return
			default:
				log.Println("Accept error:", err)
				continue
			}
		}
		go s.handleConn(conn)
	}
}
