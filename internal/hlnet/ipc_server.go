package hlnet

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"decaffeinated/pkg/net"
	"net"
)

// IPCHandler defines the function signature that will process the commands.
type IPCHandler func(cmd IPCCommand) (IPCResponse, error)// Server handles the IPC communication.

type Server struct {
	conf    IPCConfig
	handler IPCHandler
	ln      net.Listener
	done    chan struct{}
	mu      sync.Mutex
}

// NewServer creates a new instance of the IPC server.
func NewServer(conf IPCConfig, handler IPCHandler) (*Server, error) {
	if conf.Path == "" {
		return nil, errors.New("ipc path required")
	}
	if handler == nil {
		return nil, errors.New("ipc handler required")
	}
	return &Server{conf: conf, handler: handler}, nil
}

// Start opens the socket and begins listening for incoming commands.
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ln != nil {
		return errors.New("server already started")
	}

	// Using your custom lib: ListenChannels handles os.Remove(path) and net.Listen
	ln, err := dnet.ListenChannels(s.conf.Path)
	if err != nil {
		return err
	}

	s.ln = ln
	s.done = make(chan struct{})
	
	go s.serveLoop()
	return nil
}

// Stop closes the listener and shuts down the server loop.
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ln == nil {
		return nil
	}
	
	close(s.done)
	err := s.ln.Close()
	s.ln = nil
	return err
}

// serveLoop continuously accepts new connections.
func (s *Server) serveLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.done:
				return
			default:
				log.Println("IPC Accept error:", err)
				continue
			}
		}
		go s.handleConn(conn)
	}
}

// handleConn manages the lifecycle of a single connection.
func (s *Server) handleConn(c net.Conn) {
	defer c.Close()

	// Using Decoder to avoid blocking (unlike io.ReadAll)
	var cmd IPCCommand
	if err := json.NewDecoder(c).Decode(&cmd); err != nil {
		log.Printf("IPC decoding failed: %v", err)
		return
	}

	log.Printf("Executing action: %s", cmd.Action)
	
	// Execute the handler logic
	resp, err := s.handler(cmd)
	if err != nil {
		resp = IPCResponse{
			Status:  "error",
			Message: err.Error(),
		}
	}
	
	// Encode and write the response back through the same connection
	if err := json.NewEncoder(c).Encode(resp); err != nil {
		log.Printf("Failed to send IPC response: %v", err)
	}
}
