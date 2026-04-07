// Package req contains client-to-server protocol request messages for the login server.
// This file defines character creation protocol messages.
package request

import "github.com/boyism80/fm/stream"

// CreateCharacter represents a character creation request from client to login server.
// Handles new character creation with customization options for appearance and starting equipment.
//
// Flow: GameClient -> LoginServerActor -> Character validation -> Database storage -> Response
// Direction: Client -> LoginServer
// Packet Processing: Deserialize from client -> validate name/appearance -> create character
// Trigger: Player clicks "Create Character" button in character selection screen
// Side Effects:
//   - New character record created in database
//   - Character appears in player's character list
//   - Starting equipment and stats assigned
//   - Character slot consumed from player's account
//
// Error Conditions:
//   - Character name already exists
//   - Invalid character name (profanity, length, special characters)
//   - Invalid appearance combinations
//   - Account has reached maximum character limit
//   - Database connection failure
//
// Related Messages: CreateCharacterResponse (success/failure response to client)
//
// Security Notes:
//   - Character name must be validated for uniqueness and appropriateness
//   - Appearance values must be validated against allowed ranges
//   - Rate limiting should prevent character creation spam
type CreateCharacter struct {
	Name   string // Character name (3-12 characters, alphanumeric only)
	Face   uint32 // Face appearance ID (must be valid face from WZ data)
	Hair   uint32 // Hair appearance ID (must be valid hair from WZ data)
	Top    uint32 // Starting top equipment ID (must be valid beginner equipment)
	Bottom uint32 // Starting bottom equipment ID (must be valid beginner equipment)
	Shoes  uint32 // Starting shoes equipment ID (must be valid beginner equipment)
	Weapon uint32 // Starting weapon equipment ID (must be valid beginner weapon)
}

// Serialize writes the CreateCharacter message to a stream writer.
// This is typically not used for request messages (client -> server).
//
// Parameters:
//   - writer: Stream writer to serialize the message to
//
// Returns:
//   - error: Always returns nil for request messages
func (a *CreateCharacter) Serialize(writer *stream.StreamWriter) error {
	return nil
}

// Deserialize reads the CreateCharacter message from a stream reader.
// Parses the incoming client packet to extract character creation parameters.
//
// Packet Structure (in order):
//   - Name: UTF-16 string with length prefix
//   - Face: 32-bit unsigned integer (face appearance ID)
//   - Hair: 32-bit unsigned integer (hair appearance ID)
//   - Top: 32-bit unsigned integer (starting top equipment ID)
//   - Bottom: 32-bit unsigned integer (starting bottom equipment ID)
//   - Shoes: 32-bit unsigned integer (starting shoes equipment ID)
//   - Weapon: 32-bit unsigned integer (starting weapon equipment ID)
//
// Parameters:
//   - reader: Stream reader containing the client packet data
//
// Returns:
//   - error: Error if packet parsing fails or data is invalid
//
// Error Conditions:
//   - Failed to read any of the required fields
//   - Malformed packet data
//   - Unexpected end of packet
func (a *CreateCharacter) Deserialize(reader *stream.StreamReader) {
	a.Name = reader.ReadStr16()
	a.Face = reader.ReadU32()
	a.Hair = reader.ReadU32()
	a.Top = reader.ReadU32()
	a.Bottom = reader.ReadU32()
	a.Shoes = reader.ReadU32()
	a.Weapon = reader.ReadU32()

}
