package core

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/boyism80/fm/common/crypt"
	common_resp "github.com/boyism80/fm/common/protocol/resp"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/login/protocol/resp"
)

type Server struct {
	logicThreadCount   int
	listener           net.Listener
	logicThreads       []*LogicThread
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
}

// ServerConfig holds server configuration parameters
type ServerConfig struct {
	LogicThreadCount   int                                 // Number of logic threads for game processing
	Host               string                              // Server host address
	Port               int                                 // Server port number
	OnClientDisconnect func(interface{})                   // Callback when client disconnects
	ClientFactory      func(net.Conn, int) (Client, error) // Factory for creating clients
}

// NewServer creates a new server with specified thread configuration
func NewServer(config *ServerConfig) (*Server, error) {
	if config.LogicThreadCount <= 0 {
		return nil, fmt.Errorf("invalid logic thread count: %d", config.LogicThreadCount)
	}

	ctx, cancel := context.WithCancel(context.Background())

	server := &Server{
		logicThreadCount: config.LogicThreadCount,
		clients:          make(map[net.Conn]Client),
		shutdownChan:     make(chan struct{}),
		ctx:              ctx,
		cancel:           cancel,
		nextClientID:     0,
		packetHandler:    NewPacketHandler(),
		clientFactory:    config.ClientFactory,
		onClientDisconnect: func(client Client) {
			if config.OnClientDisconnect != nil {
				config.OnClientDisconnect(client)
			}
		},
	}

	return server, nil
}

// Start initializes and starts the server with configured threads
func (s *Server) Start(host string, port int) error {
	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	s.listener = listener

	// Start logic threads
	s.logicThreads = make([]*LogicThread, s.logicThreadCount)
	for i := 0; i < s.logicThreadCount; i++ {
		logicThread := &LogicThread{
			id:           i,
			server:       s,
			stopChan:     make(chan struct{}),
			taskChan:     make(chan *LogicTask, 100), // Buffer for 100 tasks
			timerManager: NewTimerManager(),
		}
		s.logicThreads[i] = logicThread

		s.wg.Add(1)
		go func(thread *LogicThread) {
			defer s.wg.Done()
			thread.run()
		}(logicThread)
	}

	// Start connection acceptor
	s.wg.Add(1)
	go s.acceptConnections()

	log.Printf("Server started with %d logic threads on %s:%d",
		s.logicThreadCount, host, port)

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

	// Cancel all timers and stop all threads
	for _, thread := range s.logicThreads {
		// Cancel all active timers
		thread.timerManager.CancelAllTimers()

		close(thread.stopChan)
	}

	// Wait for all threads to finish
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
	welcome := &common_resp.Welcome{
		SendIv: client.GetSendEncryption().IV(),
		RecvIv: client.GetRecvEncryption().IV(),
	}
	if err := client.Send(welcome, types.SEND_POLICY_RAW); err != nil {
		log.Printf("Failed to send welcome packet to %s: %v", client.GetConnection().RemoteAddr(), err)
	} else {
		log.Printf("Sent welcome packet to %s", client.GetConnection().RemoteAddr())
	}

	// Send login failed packet (for login server)
	loginFailed := &resp.LoginFailed{
		Reason: resp.LoginFailedReasonNoPopup,
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
			packetProcessed, err := s.readPacket(client)
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

			if packetProcessed {
				// Packet was submitted to logic thread for processing
				log.Printf("Packet submitted to logic thread for client %s", client.GetConnection().RemoteAddr())
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

// processPacket handles packet decryption and logic thread submission
func (s *Server) processPacket(client Client, encryptedData []byte) (bool, error) {
	// Decrypt packet data
	packetData := client.GetRecvEncryption().Decrypt(encryptedData)

	// Parse opcode (first 2 bytes)
	if len(packetData) < 2 {
		return false, fmt.Errorf("packet too short for opcode")
	}
	opcode := int(packetData[0]) | int(packetData[1])<<8
	packetData = packetData[2:] // Remove opcode from data

	// Submit packet processing to logic thread
	task := &LogicTask{
		Predicate: func() bool {
			// Check if client is still connected and valid
			return client.GetConnection() != nil
		},
		Logic: func() error {
			// Create client context
			ctx := &ClientContext{
				Client: client,
				Server: s,
				SendFunc: func(p types.Packet, policy types.SendPolicy) error {
					return client.Send(p, policy)
				},
			}

			// Handle packet in logic thread
			return s.packetHandler.Handle(ctx, opcode, packetData)
		},
		Callback: func(success bool, err error) {
			if err != nil {
				log.Printf("Error processing packet opcode 0x%02X: %v", opcode, err)
			} else {
				log.Printf("Packet processed successfully for client %s", client.GetConnection().RemoteAddr())
			}
		},
		Object:     client, // Use client for thread assignment
		MaxRetries: 3,
	}

	// Submit to appropriate logic thread
	if err := s.SubmitLogicTaskForObject(client, task); err != nil {
		log.Printf("Failed to submit packet task: %v", err)
		return false, err
	}

	return true, nil
}

// SubmitLogicTask submits a task to a logic thread (round-robin distribution)
func (s *Server) SubmitLogicTask(task *LogicTask) error {
	if task == nil {
		return fmt.Errorf("cannot submit nil task")
	}

	// Simple round-robin distribution
	staticIndex := 0 // This could be made more sophisticated
	logicThread := s.logicThreads[staticIndex%s.logicThreadCount]

	return logicThread.SubmitTask(task)
}

// SubmitLogicTaskToThread submits a task to a specific logic thread
func (s *Server) SubmitLogicTaskToThread(threadIndex int, task *LogicTask) error {
	if task == nil {
		return fmt.Errorf("cannot submit nil task")
	}

	if threadIndex < 0 || threadIndex >= s.logicThreadCount {
		return fmt.Errorf("invalid thread index: %d (valid range: 0-%d)", threadIndex, s.logicThreadCount-1)
	}

	return s.logicThreads[threadIndex].SubmitTask(task)
}

// SubmitLogicTaskForObject submits a task to the appropriate logic thread for a ThreadAssignable object
func (s *Server) SubmitLogicTaskForObject(obj ThreadAssignable, task *LogicTask) error {
	if obj == nil {
		return fmt.Errorf("cannot submit task for nil object")
	}
	if task == nil {
		return fmt.Errorf("cannot submit nil task")
	}

	// Set default values for thread reassignment
	if task.Object == nil {
		task.Object = obj
	}
	if task.MaxRetries == 0 {
		task.MaxRetries = 3 // Default max retries
	}
	task.retryCount = 0 // Reset retry count

	// Calculate thread index using object's hash
	threadIndex := obj.GetThreadHash() % s.logicThreadCount
	return s.SubmitLogicTaskToThread(threadIndex, task)
}

// GetThreadForObject returns the logic thread index for a ThreadAssignable object
func (s *Server) GetThreadForObject(obj ThreadAssignable) (int, error) {
	if obj == nil {
		return -1, fmt.Errorf("cannot get thread for nil object")
	}

	threadIndex := obj.GetThreadHash() % s.logicThreadCount
	return threadIndex, nil
}

// GetLogicThread returns the LogicThread instance for a ThreadAssignable object
func (s *Server) GetLogicThread(obj ThreadAssignable) (*LogicThread, error) {
	if obj == nil {
		return nil, fmt.Errorf("cannot get thread for nil object")
	}

	threadIndex := obj.GetThreadHash() % s.logicThreadCount
	if threadIndex < 0 || threadIndex >= s.logicThreadCount {
		return nil, fmt.Errorf("invalid thread index: %d (valid range: 0-%d)", threadIndex, s.logicThreadCount-1)
	}

	return s.logicThreads[threadIndex], nil
}

// SetTimer sets a repeating timer on all logic threads
// This is useful for server-wide periodic tasks
func (s *Server) SetTimer(interval time.Duration, logic func() error, callback func(bool, error)) []*RepeatingTimer {
	var timers []*RepeatingTimer

	for _, logicThread := range s.logicThreads {
		timer := logicThread.SetRepeatingTimer(interval, logic, callback)
		timers = append(timers, timer)
	}

	return timers
}

// CancelTimer cancels a repeating timer by ID on all logic threads
func (s *Server) CancelTimer(timerID uint64) error {
	var lastError error
	successCount := 0

	for _, logicThread := range s.logicThreads {
		if err := logicThread.CancelRepeatingTimer(timerID); err == nil {
			successCount++
		} else {
			lastError = err
		}
	}

	if successCount == 0 {
		if lastError != nil {
			return fmt.Errorf("timer %d not found on any logic thread: %w", timerID, lastError)
		}
		return fmt.Errorf("timer %d not found on any logic thread", timerID)
	}

	return nil
}

// GetRepeatingTimerCount returns the total number of repeating timers across all logic threads
func (s *Server) GetRepeatingTimerCount() int {
	total := 0
	for _, logicThread := range s.logicThreads {
		total += logicThread.GetRepeatingTimerCount()
	}
	return total
}

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

	// Get timer statistics for each logic thread
	timerStats := make(map[string]int)
	repeatingTimerStats := make(map[string]int)
	totalTimers := 0
	totalRepeatingTimers := 0
	for i, thread := range s.logicThreads {
		timerCount := thread.GetTimerCount()
		repeatingTimerCount := thread.GetRepeatingTimerCount()
		timerStats[fmt.Sprintf("thread_%d_timers", i)] = timerCount
		repeatingTimerStats[fmt.Sprintf("thread_%d_repeating_timers", i)] = repeatingTimerCount
		totalTimers += timerCount
		totalRepeatingTimers += repeatingTimerCount
	}

	return map[string]interface{}{
		"logic_thread_count":     s.logicThreadCount,
		"client_count":           clientCount,
		"listening":              s.listener != nil,
		"packet_handlers":        s.packetHandler.GetHandlerCount(),
		"total_timers":           totalTimers,
		"total_repeating_timers": totalRepeatingTimers,
		"timer_stats":            timerStats,
		"repeating_timer_stats":  repeatingTimerStats,
	}
}

// SetOnClientDisconnect sets the callback function for client disconnection
func (s *Server) SetOnClientDisconnect(callback func(interface{})) {
	s.onClientDisconnect = func(client Client) {
		callback(client)
	}
}
