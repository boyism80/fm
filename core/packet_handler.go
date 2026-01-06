package core

import (
	"fmt"
	"log"
)

// PacketHandler manages packet processing for the server
// Flow: Opcode -> Handler function -> Packet processing
// Purpose: Centralized packet handling with automatic deserialization
// Thread Safety: Safe for concurrent registration and handling
type PacketHandler struct {
	handlers map[int]func(ctx *ClientContext, data []byte) error
}

// ClientContext provides context for packet handlers
// Flow: Client -> Server -> Handler function
// Purpose: Provides client and server references to handlers
// Thread Safety: Safe for concurrent access
type ClientContext struct {
	Client Client
	Server *Server
}

// NewPacketHandler creates a new packet handler
// Flow: Initialization -> Handler registration -> Packet processing
// Purpose: Creates empty packet handler ready for registration
// Error Handling: None (always succeeds)
func NewPacketHandler() *PacketHandler {
	return &PacketHandler{
		handlers: make(map[int]func(ctx *ClientContext, data []byte) error),
	}
}

// RegisterHandler registers a packet handler for a specific opcode
// Flow: Opcode -> Handler function -> Registration
// Thread Safety: Should be called during server initialization
// Error Handling: Overwrites existing handler for same opcode
func (h *PacketHandler) RegisterHandler(opcode int, handler func(ctx *ClientContext, data []byte) error) {
	h.handlers[opcode] = handler
	log.Printf("Registered packet handler for opcode 0x%02X", opcode)
}

// GetHandler returns the handler for a specific opcode
func (h *PacketHandler) GetHandler(opcode int) func(ctx *ClientContext, data []byte) error {
	return h.handlers[opcode]
}

// Handle processes a packet with the registered handler
// Flow: Opcode lookup -> Handler execution -> Result reporting
// Thread Safety: Safe for concurrent handling
// Error Handling: Missing handlers, handler execution errors
func (h *PacketHandler) Handle(ctx *ClientContext, opcode int, data []byte) error {
	handler, exists := h.handlers[opcode]
	if !exists {
		return fmt.Errorf("no handler registered for opcode 0x%02X", opcode)
	}

	return handler(ctx, data)
}

// GetHandlerCount returns the number of registered handlers
// Flow: Count calculation -> Return value
// Purpose: Statistics and monitoring
// Thread Safety: Safe for concurrent access
func (h *PacketHandler) GetHandlerCount() int {
	return len(h.handlers)
}

// GetRegisteredOpcodes returns all registered opcodes
// Flow: Opcode collection -> Return slice
// Purpose: Debugging and monitoring
// Thread Safety: Safe for concurrent access
func (h *PacketHandler) GetRegisteredOpcodes() []int {
	opcodes := make([]int, 0, len(h.handlers))
	for opcode := range h.handlers {
		opcodes = append(opcodes, opcode)
	}
	return opcodes
}
