// Package resp contains server-to-client protocol response messages for the login server.
// This file defines character creation response protocol messages.
package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/entity"
)

// CreateCharacter represents a character creation response.
// Flow: LoginServerActor -> GameClient (after processing creation request)
// Packet ID: 0x06
// Purpose: Inform client whether character creation succeeded and provide character data
// Side Effects:
//   - On success: Character appears in character selection list
//   - On failure: Client displays error message to user
//   - Character slot count updated in client
//
// Error Conditions:
//   - Character name already exists (Success = false)
//   - Invalid character data (Success = false)
//   - Database error (Success = false)
//
// Related Messages: CreateCharacter (request), CharacterList (updated list after creation)
//
// Security Notes:
//   - Character data must be properly validated before creation
//   - Success flag must accurately reflect creation result
//   - Character overview data should match what was actually created
type CreateCharacter struct {
	Character *entity.Character // Created character data (nil if failed)
	Success   bool              // Whether creation was successful
}

// Opcode returns the packet opcode for CreateCharacter
func (a *CreateCharacter) Opcode() uint16 {
	return 0x06
}

// Serialize writes the character creation result packet to client.
// Constructs the character creation result packet to send to the client.
//
// Packet Structure (in order):
//   - Packet ID: 0x06 (16-bit)
//   - Error flag: !Success (8-bit boolean, true indicates error)
//   - Character overview: Serialized character data for display in character list
//
// The character overview contains basic character information needed for the
// character selection screen, including name, level, job, appearance, etc.
//
// Parameters:
//   - writer: Stream writer to serialize the response to
//
// Returns:
//   - error: Error if serialization fails (typically nil)
//
// Implementation Notes:
//   - Error flag is inverted (!Success) to match client expectations
//   - Character.SerializeOverview() handles the character data serialization
//   - If Success is false, Character may be nil or contain partial data
func (a *CreateCharacter) Serialize(writer *stream.StreamWriter) error {
	writer.WriteBoolean(!a.Success)
	a.Character.SerializeOverview(writer)
	return nil
}

// Deserialize reads from stream (not used for response messages).
// This is typically not used for response messages (server -> client).
//
// Parameters:
//   - reader: Stream reader containing packet data
//
// Returns:
//   - error: Always returns nil for response messages
func (a *CreateCharacter) Deserialize(reader *stream.StreamReader) error {
	return nil
}
