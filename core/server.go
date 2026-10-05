package core

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
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
)

var packetLogEnabled atomic.Bool

func SetPacketLogEnabled(b bool) {
	packetLogEnabled.Store(b)
}

func GetPacketLogEnabled() bool {
	return packetLogEnabled.Load()
}

type ServerCore struct {
	listener           net.Listener
	clients            map[net.Conn]Client
	clientsMutex       sync.RWMutex
	shutdownChan       chan struct{}
	wg                 sync.WaitGroup
	ctx                context.Context
	cancel             context.CancelFunc
	nextClientID       int
	clientIDMutex      sync.Mutex
	packetHandler      *PacketHandler
	onClientDisconnect func(Client)
	clientFactory      func(net.Conn, int) (Client, error)
	config             *ServerConfig
	rootContext        *actor.RootContext
	nilActorPID        *actor.PID
}

type ServerConfig struct {
	Host               string
	Port               int
	OnClientConnect    func(Client)
	OnClientDisconnect func(Client)
	ClientFactory      func(net.Conn, int) (Client, error)
}

func NewServer(config *ServerConfig) (*ServerCore, error) {
	ctx, cancel := context.WithCancel(context.Background())

	server := &ServerCore{
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

func (s *ServerCore) SetRootContext(root *actor.RootContext) {
	s.rootContext = root
}

func (s *ServerCore) SetNilActorPID(pid *actor.PID) {
	s.nilActorPID = pid
}

func (s *ServerCore) GetRootContext() *actor.RootContext {
	return s.rootContext
}

func (s *ServerCore) Start(host string, port int) error {

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	s.listener = listener

	s.wg.Add(1)
	go s.acceptConnections()

	log.Printf("Server started on %s:%d", host, port)

	return nil
}

func (s *ServerCore) Stop() error {
	log.Println("Shutting down server...")

	close(s.shutdownChan)
	s.cancel()

	if s.listener != nil {
		s.listener.Close()
	}

	s.clientsMutex.Lock()
	for conn, client := range s.clients {
		client.Close()
		delete(s.clients, conn)
	}
	s.clientsMutex.Unlock()

	s.wg.Wait()

	log.Println("Server stopped successfully")
	return nil
}

func (s *ServerCore) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdownChan:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if s.ctx.Err() != nil {

					return
				}
				log.Printf("Failed to accept connection: %v", err)
				continue
			}

			s.clientIDMutex.Lock()
			clientID := s.nextClientID
			s.nextClientID++
			s.clientIDMutex.Unlock()

			client, err := s.createClient(conn, clientID)
			if err != nil {
				log.Printf("Failed to create client for %s: %v", conn.RemoteAddr(), err)
				conn.Close()
				continue
			}

			s.clientsMutex.Lock()
			s.clients[conn] = client
			s.clientsMutex.Unlock()

			if s.config.OnClientConnect != nil {
				s.config.OnClientConnect(client)
			}

			s.wg.Add(1)
			go s.handleClient(client)

			log.Printf("New client connected: %s", conn.RemoteAddr())
		}
	}
}

func (s *ServerCore) createClient(conn net.Conn, clientID int) (Client, error) {
	if s.clientFactory != nil {
		return s.clientFactory(conn, clientID)
	}
	return nil, fmt.Errorf("no client factory provided")
}

func (s *ServerCore) handleClient(client Client) {
	defer func() {

		if s.onClientDisconnect != nil {
			s.onClientDisconnect(client)
		}

		s.clientsMutex.Lock()
		delete(s.clients, client.GetConnection())
		s.clientsMutex.Unlock()
		client.Close()
		s.wg.Done()
		log.Printf("Client disconnected: %s", client.GetConnection().RemoteAddr())
	}()

	welcome := &response.Welcome{
		SendIv: client.GetSendEncryption().IV(),
		RecvIv: client.GetRecvEncryption().IV(),
	}
	if err := client.Send(welcome, types.SEND_POLICY_RAW); err != nil {
		log.Printf("Failed to send welcome packet to %s: %v", client.GetConnection().RemoteAddr(), err)
	} else {
		log.Printf("Sent welcome packet to %s", client.GetConnection().RemoteAddr())
	}

	loginFailed := &response.LoginFailed{
		Reason: response.LoginFailedReasonNoPopup,
	}
	if err := client.Send(loginFailed, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send login failed packet to %s: %v", client.GetConnection().RemoteAddr(), err)
	} else {
		log.Printf("Sent login failed packet to %s", client.GetConnection().RemoteAddr())
	}

	reader := bufio.NewReader(client.GetConnection())
	for {
		select {
		case <-s.shutdownChan:
			return
		default:

			err := s.readPacket(client, reader)
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Printf("Client disconnected normally: %s", client.GetConnection().RemoteAddr())
					return
				}
				if errors.Is(err, net.ErrClosed) {
					log.Printf("Client disconnected (connection closed): %s", client.GetConnection().RemoteAddr())
					return
				}

				if strings.Contains(err.Error(), "wsarecv") || strings.Contains(err.Error(), "connection was aborted") {
					log.Printf("Client connection aborted: %s", client.GetConnection().RemoteAddr())
					return
				}

				log.Printf("Error reading packet from %s: %v", client.GetConnection().RemoteAddr(), err)
				return
			}
		}
	}
}

func (s *ServerCore) readPacket(client Client, reader *bufio.Reader) error {

	headerBytes := make([]byte, 4)
	if _, err := io.ReadFull(reader, headerBytes); err != nil {
		return err
	}

	if !client.GetRecvEncryption().CheckPacketHeader(headerBytes) {
		return fmt.Errorf("invalid packet header")
	}

	packetLength, err := crypt.GetPacketLength(headerBytes)
	if err != nil {
		return fmt.Errorf("failed to get packet length: %w", err)
	}

	if packetLength <= 0 || packetLength > 8192 {
		return fmt.Errorf("invalid packet length: %d", packetLength)
	}

	encryptedData := make([]byte, packetLength)
	if _, err := io.ReadFull(reader, encryptedData); err != nil {
		return err
	}

	return s.processPacket(client, encryptedData)
}

func (s *ServerCore) processPacket(client Client, encryptedData []byte) error {

	packetData := client.GetRecvEncryption().Decrypt(encryptedData)

	if len(packetData) < 2 {
		return fmt.Errorf("packet too short for opcode")
	}
	opcode := int(packetData[0]) | int(packetData[1])<<8
	packetData = packetData[2:]

	logicActorPID := s.PacketActorPID(client)
	if logicActorPID == nil {
		log.Printf("No LogicActor PID for client and nil LogicActor PID not set")
		return fmt.Errorf("no LogicActor PID for client and nil LogicActor PID not set")
	}

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
		return fmt.Errorf("rootContext not set")
	}

	return nil
}

// PacketActorPID is the actor that handles the client's packets now: its map's actor, or the nil actor between maps.
func (s *ServerCore) PacketActorPID(client Client) *actor.PID {
	if logicClient, ok := client.(interface{ GetLogicActorPID() *actor.PID }); ok {
		if pid := logicClient.GetLogicActorPID(); pid != nil {
			return pid
		}
	}
	return s.nilActorPID
}

func (s *ServerCore) GetPacketHandler() *PacketHandler {
	return s.packetHandler
}

func (s *ServerCore) GetStats() map[string]interface{} {
	s.clientsMutex.RLock()
	clientCount := len(s.clients)
	s.clientsMutex.RUnlock()

	return map[string]interface{}{
		"client_count":    clientCount,
		"listening":       s.listener != nil,
		"packet_handlers": s.packetHandler.GetHandlerCount(),
	}
}

func (s *ServerCore) Clients() []Client {
	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()
	clients := make([]Client, 0, len(s.clients))
	for _, c := range s.clients {
		clients = append(clients, c)
	}
	return clients
}

func (s *ServerCore) SetOnClientDisconnect(callback func(Client)) {
	s.onClientDisconnect = callback
}
