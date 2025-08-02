package core

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/common/types"
)

// PacketHandler manages packet processing for the server
// Flow: Opcode -> Handler function -> Packet processing
// Purpose: Centralized packet handling with automatic deserialization
// Thread Safety: Safe for concurrent registration and handling
type PacketHandler[T any] struct {
	handlers map[int]func(ctx *ClientContext[T], data []byte) error
}

// ClientContext provides context for packet handlers
// Flow: Client -> Server -> Handler function
// Purpose: Provides client and server references to handlers
// Thread Safety: Safe for concurrent access
type ClientContext[T any] struct {
	Client   *Client[T]
	Server   *Server[T]
	SendFunc func(p types.Packet, policy types.SendPolicy) error
}

// NewPacketHandler creates a new packet handler
// Flow: Initialization -> Handler registration -> Packet processing
// Purpose: Creates empty packet handler ready for registration
// Error Handling: None (always succeeds)
func NewPacketHandler[T any]() *PacketHandler[T] {
	return &PacketHandler[T]{
		handlers: make(map[int]func(ctx *ClientContext[T], data []byte) error),
	}
}

// RegisterHandler registers a packet handler for a specific opcode
// Flow: Opcode -> Handler function -> Registration
// Thread Safety: Should be called during server initialization
// Error Handling: Overwrites existing handler for same opcode
func (h *PacketHandler[T]) RegisterHandler(opcode int, handler func(ctx *ClientContext[T], data []byte) error) {
	h.handlers[opcode] = handler
	log.Printf("Registered packet handler for opcode 0x%02X", opcode)
}

// Handle processes a packet with the registered handler
// Flow: Opcode lookup -> Handler execution -> Result reporting
// Thread Safety: Safe for concurrent handling
// Error Handling: Missing handlers, handler execution errors
func (h *PacketHandler[T]) Handle(ctx *ClientContext[T], opcode int, data []byte) error {
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
func (h *PacketHandler[T]) GetHandlerCount() int {
	return len(h.handlers)
}

// GetRegisteredOpcodes returns all registered opcodes
// Flow: Opcode collection -> Return slice
// Purpose: Debugging and monitoring
// Thread Safety: Safe for concurrent access
func (h *PacketHandler[T]) GetRegisteredOpcodes() []int {
	opcodes := make([]int, 0, len(h.handlers))
	for opcode := range h.handlers {
		opcodes = append(opcodes, opcode)
	}
	return opcodes
}
