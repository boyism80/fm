package req

import (
	"github.com/boyism80/fm/common/stream"
)

// Revive represents a player revival request packet
// Direction: Client -> GameServer
// Packet ID: 0x2B (ChangeMap)
// Response: Warp response with revival map
// Security: Validates player state and wheel of fortune item
type Revive struct {
	UnknownFlag uint8  // Unknown flag byte
	TargetID    int32  // Target map ID (-1 for normal portal, positive for revival)
	PortalName  string // Portal name for movement
	SkipByte    uint8  // Skipped byte
	UseWheel    bool   // Whether to use wheel of fortune for revival
}

// Serialize writes the revive packet to the stream
// Note: This is a client->server packet, so serialization is not typically used
func (p *Revive) Serialize(writer *stream.StreamWriter) error {
	return nil
}

// Deserialize reads the revive packet from the stream
// Packet structure:
// 1. UnknownFlag (1 byte) - Unknown flag
// 2. TargetID (4 bytes) - Target map ID (-1 for portal, positive for revival)
// 3. PortalName (string) - Portal name for movement
// 4. SkipByte (1 byte) - Skipped byte
// 5. UseWheel (1 byte) - Wheel of fortune usage flag
func (p *Revive) Deserialize(reader *stream.StreamReader) error {
	p.UnknownFlag, _ = reader.ReadU8()
	p.TargetID, _ = reader.Read32()
	p.PortalName, _ = reader.ReadStr16()
	reader.Skip(1) // Skip 1 byte
	useWheelByte, _ := reader.ReadU8()
	p.UseWheel = useWheelByte > 0
	return nil
} 