// Package msg contains common message types used across the MapleStory private server.
// This file defines network connectivity and health check messages.
package msg

// Ping represents a network connectivity test message.
// Flow: Any Actor -> Target Actor
// Used for connection health checks and keepalive mechanisms.
//
// Trigger: Network health checks, connection validation, keepalive mechanisms
// Purpose: Ensures network connectivity and measures round-trip time
// Side Effects: May trigger connection cleanup if no response received
// Related Messages: Pong (response message, if implemented)
//
// Usage Examples:
//   - Client-Server keepalive
//   - Inter-server connectivity checks
//   - Network latency measurement
type Ping struct {
	// No fields required - presence of message indicates ping request
}
