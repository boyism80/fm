package core

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/crypt"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
	// lua "github.com/yuin/gopher-lua" // Commented out: LogicThread removed
)

// packetLogEnabled controls whether each received client packet is logged (opcode + payload). Default false.
var packetLogEnabled atomic.Bool

// SetPacketLogEnabled sets whether to log received client packets. Call from Lua builtin or elsewhere to toggle.
func SetPacketLogEnabled(b bool) {
	packetLogEnabled.Store(b)
}

// GetPacketLogEnabled returns whether received client packets are currently logged.
func GetPacketLogEnabled() bool {
	return packetLogEnabled.Load()
}

type Server struct {
	listener           net.Listener
	clients            map[net.Conn]Client
	clientsMutex       sync.RWMutex
	shutdownChan       chan struct{}
	wg                 sync.WaitGroup
	ctx                context.Context
	cancel             context.CancelFunc
	nextClientID       int // Counter for generating unique client IDs
	clientIDMutex      sync.Mutex
	packetHandler      *PacketHandler                      // Central packet handler for the server
	onClientDisconnect func(Client)                        // Callback when client disconnects
	clientFactory      func(net.Conn, int) (Client, error) // Factory for creating clients
	config             *ServerConfig                       // Server configuration
	rootContext        *actor.RootContext                  // RootContext for sending messages to actors
	nilMapActorPID     *actor.PID                          // PID for nil MapActor (for characters before map assignment)
}

// ServerConfig holds server configuration parameters
type ServerConfig struct {
	Host               string                              // Server host address
	Port               int                                 // Server port number
	OnClientConnect    func(Client)                        // Callback when client connects
	OnClientDisconnect func(Client)                        // Callback when client disconnects
	ClientFactory      func(net.Conn, int) (Client, error) // Factory for creating clients
	// LogicThreadInit    func(*LogicThread)                  // Commented out: LogicThread removed
}

// NewServer creates a new server
func NewServer(config *ServerConfig) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	server := &Server{
		clients:       make(map[net.Conn]Client),
		shutdownChan:  make(chan struct{}),
		ctx:           ctx,
		cancel:        cancel,
		nextClientID:  0,
		packetHandler: NewPacketHandler(),
		clientFactory: config.ClientFactory,
		onClientDisconnect: func(client Client) {
			if config.OnClientDisconnect != nil {
				config.OnClientDisconnect(client)
			}
		},
		config: config,
	}

	return server, nil
}

// SetRootContext sets the RootContext for sending messages to actors
func (s *Server) SetRootContext(root *actor.RootContext) {
	s.rootContext = root
}

// SetNilMapActorPID sets the PID for nil MapActor
func (s *Server) SetNilMapActorPID(pid *actor.PID) {
	s.nilMapActorPID = pid
}

// GetRootContext returns the RootContext for sending messages to actors
func (s *Server) GetRootContext() *actor.RootContext {
	return s.rootContext
}

// Start initializes and starts the server with configured threads
func (s *Server) Start(host string, port int) error {
	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	s.listener = listener

	// Start connection acceptor
	s.wg.Add(1)
	go s.acceptConnections()

	log.Printf("Server started on %s:%d", host, port)

	return nil
}

// Stop gracefully shuts down the server and all threads
func (s *Server) Stop() error {
	log.Println("Shutting down server...")

	// Signal shutdown
	close(s.shutdownChan)
	s.cancel()

	// Stop accepting new connections
	if s.listener != nil {
		s.listener.Close()
	}

	// Close all client connections
	s.clientsMutex.Lock()
	for conn, client := range s.clients {
		client.GetConnection().Close()
		delete(s.clients, conn)
	}
	s.clientsMutex.Unlock()

	// Wait for all goroutines to finish
	s.wg.Wait()

	log.Println("Server stopped successfully")
	return nil
}

// acceptConnections accepts new client connections and starts goroutines for each
func (s *Server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdownChan:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if s.ctx.Err() != nil {
					// Server is shutting down
					return
				}
				log.Printf("Failed to accept connection: %v", err)
				continue
			}

			// Generate unique client ID
			s.clientIDMutex.Lock()
			clientID := s.nextClientID
			s.nextClientID++
			s.clientIDMutex.Unlock()

			// Create client using factory
			client, err := s.createClient(conn, clientID)
			if err != nil {
				log.Printf("Failed to create client for %s: %v", conn.RemoteAddr(), err)
				conn.Close()
				continue
			}

			// Add to server client list
			s.clientsMutex.Lock()
			s.clients[conn] = client
			s.clientsMutex.Unlock()

			// Call OnClientConnect callback if set
			if s.config.OnClientConnect != nil {
				s.config.OnClientConnect(client)
			}

			// Start goroutine for this client
			s.wg.Add(1)
			go s.handleClient(client)

			log.Printf("New client connected: %s", conn.RemoteAddr())
		}
	}
}

// createClient creates a new client using the configured factory
func (s *Server) createClient(conn net.Conn, clientID int) (Client, error) {
	if s.clientFactory != nil {
		return s.clientFactory(conn, clientID)
	}
	return nil, fmt.Errorf("no client factory provided")
}

// handleClient handles a single client connection in its own goroutine
func (s *Server) handleClient(client Client) {
	defer func() {
		// Call disconnect callback if set
		if s.onClientDisconnect != nil {
			s.onClientDisconnect(client)
		}

		// Clean up client
		s.clientsMutex.Lock()
		delete(s.clients, client.GetConnection())
		s.clientsMutex.Unlock()
		client.GetConnection().Close()
		s.wg.Done()
		log.Printf("Client disconnected: %s", client.GetConnection().RemoteAddr())
	}()

	// Send welcome packet on connection (for login server)
	welcome := &response.Welcome{
		SendIv: client.GetSendEncryption().IV(),
		RecvIv: client.GetRecvEncryption().IV(),
	}
	if err := client.Send(welcome, types.SEND_POLICY_RAW); err != nil {
		log.Printf("Failed to send welcome packet to %s: %v", client.GetConnection().RemoteAddr(), err)
	} else {
		log.Printf("Sent welcome packet to %s", client.GetConnection().RemoteAddr())
	}

	// Send login failed packet (for login server)
	loginFailed := &response.LoginFailed{
		Reason: response.LoginFailedReasonNoPopup,
	}
	if err := client.Send(loginFailed, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send login failed packet to %s: %v", client.GetConnection().RemoteAddr(), err)
	} else {
		log.Printf("Sent login failed packet to %s", client.GetConnection().RemoteAddr())
	}

	for {
		select {
		case <-s.shutdownChan:
			return
		default:
			// Try to read packet
			_, err := s.readPacket(client)
			if err != nil {
				if err.Error() == "EOF" {
					// Client disconnected normally
					log.Printf("Client disconnected normally: %s", client.GetConnection().RemoteAddr())
					return
				}

				// Check for connection abort errors (common on Windows)
				if strings.Contains(err.Error(), "wsarecv") || strings.Contains(err.Error(), "connection was aborted") {
					log.Printf("Client connection aborted: %s", client.GetConnection().RemoteAddr())
					return
				}

				log.Printf("Error reading packet from %s: %v", client.GetConnection().RemoteAddr(), err)
				continue
			}
		}
	}
}

// readPacket attempts to read a packet from the client connection
func (s *Server) readPacket(client Client) (bool, error) {
	// Read packet header (4 bytes)
	headerBytes := make([]byte, 4)
	if _, err := client.GetConnection().Read(headerBytes); err != nil {
		return false, err
	}

	// Check packet header validity
	if !client.GetRecvEncryption().CheckPacketHeader(headerBytes) {
		return false, fmt.Errorf("invalid packet header")
	}

	// Get packet length from header
	packetLength, err := crypt.GetPacketLength(headerBytes)
	if err != nil {
		return false, fmt.Errorf("failed to get packet length: %w", err)
	}

	if packetLength <= 0 || packetLength > 8192 {
		return false, fmt.Errorf("invalid packet length: %d", packetLength)
	}

	// Read encrypted packet data
	encryptedData := make([]byte, packetLength)
	if _, err := client.GetConnection().Read(encryptedData); err != nil {
		return false, err
	}

	// Process the packet (handles decryption and logic thread submission)
	return s.processPacket(client, encryptedData)
}

// processPacket handles packet decryption and actor submission
func (s *Server) processPacket(client Client, encryptedData []byte) (bool, error) {
	// Decrypt packet data
	packetData := client.GetRecvEncryption().Decrypt(encryptedData)

	// Parse opcode (first 2 bytes)
	if len(packetData) < 2 {
		return false, fmt.Errorf("packet too short for opcode")
	}
	opcode := int(packetData[0]) | int(packetData[1])<<8
	packetData = packetData[2:] // Remove opcode from data

	// Get LogicActor PID from client when supported
	var logicActorPID *actor.PID
	if logicClient, ok := client.(interface{ GetLogicActorPID() *actor.PID }); ok {
		logicActorPID = logicClient.GetLogicActorPID()
	}
	if logicActorPID == nil {
		// Use nil MapActor if PID is nil
		if s.nilMapActorPID == nil {
			log.Printf("No LogicActor PID for client and nil MapActor PID not set")
			return false, fmt.Errorf("no LogicActor PID for client and nil MapActor PID not set")
		}
		logicActorPID = s.nilMapActorPID
	}

	// Send HandlePacket message to LogicActor
	msg := &c_actor.HandlePacket{
		Opcode:        opcode,
		Data:          packetData,
		Client:        client,
		LogicActorPID: logicActorPID,
	}

	if s.rootContext != nil {
		s.rootContext.Send(logicActorPID, msg)
	} else {
		log.Printf("RootContext not set, cannot send packet to actor")
		return false, fmt.Errorf("rootContext not set")
	}

	return true, nil
}

// LogicThread related methods removed - using Actor model instead

// RegisterPacketHandler registers a packet handler for a specific opcode
func (s *Server) RegisterPacketHandler(opcode int, handler func(ctx *ClientContext, data []byte) error) {
	s.packetHandler.RegisterHandler(opcode, handler)
}

// GetPacketHandler returns the server's packet handler for external access
func (s *Server) GetPacketHandler() *PacketHandler {
	return s.packetHandler
}

// GetStats returns server statistics for monitoring
func (s *Server) GetStats() map[string]interface{} {
	s.clientsMutex.RLock()
	clientCount := len(s.clients)
	s.clientsMutex.RUnlock()

	return map[string]interface{}{
		"client_count":    clientCount,
		"listening":       s.listener != nil,
		"packet_handlers": s.packetHandler.GetHandlerCount(),
	}
}

// SetOnClientDisconnect sets the callback function for client disconnection
func (s *Server) SetOnClientDisconnect(callback func(Client)) {
	s.onClientDisconnect = callback
}
