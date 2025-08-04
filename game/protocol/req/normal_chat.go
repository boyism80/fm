// Package req contains client-to-server protocol request message definitions.
// This file defines chat-related protocol messages for player communication.
package req

import (
	"github.com/boyism80/fm/core/stream"
)

// NormalChat represents a normal chat message from client to server.
// Flow: GameClient -> GameClientActor -> MapActor -> broadcast to nearby players
// Direction: Client -> GameServer
// Packet Processing: Deserialize from client packet -> validate -> broadcast
// Trigger: Player types and sends a chat message
// Side Effects:
//   - Message broadcasted to players in same map
//   - May be recorded in chat history (unless DontRecordHistory is true)
//   - Triggers anti-spam and profanity filters
//
// Error Conditions:
//   - Message too long (exceeds character limit)
//   - Player is muted or chat-banned
//   - Invalid characters or formatting
//
// Related Messages: ChatBroadcast (response to other players)
//
// Security Notes:
//   - Message content must be validated for length and content
//   - Should be filtered for profanity and spam
//   - Rate limiting should be applied per player
type NormalChat struct {
	Message           string // Chat message content
	DontRecordHistory bool   // Whether to exclude from chat history
}

// Serialize writes the message to stream (not used for request messages).
func (m *NormalChat) Serialize(writer *stream.StreamWriter) error {
	return nil
}

// Deserialize reads the chat message from client packet.
func (m *NormalChat) Deserialize(reader *stream.StreamReader) error {
	message, err := reader.ReadStr16()
	if err != nil {
		return err
	}

	show, err := reader.ReadBool()
	if err != nil {
		return err
	}

	m.Message = message
	m.DontRecordHistory = show
	return nil
}
