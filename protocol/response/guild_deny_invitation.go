package response

import "github.com/boyism80/fm/stream"

type GuildDenyInvitation struct {
	Code         uint8
	DeclinerName string
}

func (p *GuildDenyInvitation) Opcode() uint16 {
	return guildOperationOpcode()
}

func (p *GuildDenyInvitation) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(p.Code)
	w.WriteStr16(p.DeclinerName)
	return nil
}

func (p *GuildDenyInvitation) Deserialize(*stream.StreamReader) {}
