// Package resp contains server-to-client protocol response messages for the login server.
// This file defines authentication-related protocol response messages.
package response

import (
	"github.com/boyism80/fm/stream"
)

// Authenticate represents a successful authentication response.
// Flow: LoginServerActor -> GameClient (after successful login)
// Packet ID: 0x00
// Purpose: Inform client of successful login and provide account details
// Side Effects:
//   - Client proceeds to character selection screen
//   - Account session established
//   - Chat restrictions applied if account is chat-banned
//
// Error Conditions: None (this is a success response)
// Related Messages: Login (request), CharacterList (next step after authentication)
//
// Security Notes:
//   - Contains sensitive account information
//   - Should only be sent after proper authentication
//   - Chat ban information must be accurately reflected
type Authenticate struct {
	AccountId     uint32 // Account ID in database
	Gender        uint8  // Account gender (0=male, 1=female, 2=unset)
	Admin         bool   // Administrator privileges
	AccountName   string // Account username
	IsChatBlocked bool   // Chat ban status
	ChatBlockTime uint64 // Chat ban expiration (FILETIME)
}

// Opcode returns the packet opcode for Authenticate
func (a *Authenticate) Opcode() uint16 {
	return 0x00
}

// Serialize writes the authentication success packet to client.
// Constructs the authentication success packet to send to the client.
//
// Packet Structure (in order):
//   - Packet ID: 0x00 (16-bit)
//   - Status: 0 (8-bit, indicates success)
//   - Account ID: 32-bit unsigned integer
//   - Gender: 8-bit unsigned integer
//   - Admin flag: Boolean
//   - Unknown: 0 (8-bit padding)
//   - Account name: UTF-16 string with length prefix
//   - Unknown fields: Various padding/unused fields
//   - Chat blocked flag: Boolean
//   - Chat block time: 64-bit unsigned integer (FILETIME)
//   - Additional empty strings: UTF-16 strings (likely for future use)
//
// Parameters:
//   - writer: Stream writer to serialize the response to
//
// Returns:
//   - error: Error if serialization fails (typically nil)
func (a *Authenticate) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0)
	writer.WriteU32(a.AccountId)
	writer.WriteU8(a.Gender)
	writer.WriteBoolean(a.Admin)
	writer.WriteU8(0)
	writer.WriteStr16(a.AccountName)
	writer.WriteU32(0)
	writer.WriteU8(0)
	writer.WriteU8(0)
	writer.WriteBoolean(a.IsChatBlocked)
	writer.WriteU64(a.ChatBlockTime)
	writer.WriteStr16("")
	writer.WriteStr16("")
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
func (a *Authenticate) Deserialize(reader *stream.StreamReader) {
}
