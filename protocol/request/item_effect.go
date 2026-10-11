package request

import "github.com/boyism80/fm/stream"

type ItemEffect struct{ ItemID uint32 }

func (*ItemEffect) Opcode() byte                             { return 0x23 }
func (p *ItemEffect) Serialize(w *stream.StreamWriter) error { w.WriteU32(p.ItemID); return nil }
func (p *ItemEffect) Deserialize(r *stream.StreamReader)     { p.ItemID = r.ReadU32() }
