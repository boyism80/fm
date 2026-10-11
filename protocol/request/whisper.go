package request

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type Whisper struct {
	Mode    internal.WhisperMode
	Target  string
	Message string
}

func (*Whisper) Opcode() byte { return 0x63 }

func (p *Whisper) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Mode))
	w.WriteStr16(p.Target)
	if p.Mode == internal.WhisperMode_WHISPER_MODE_SEND {
		w.WriteStr16(p.Message)
	}
	return nil
}

func (p *Whisper) Deserialize(r *stream.StreamReader) {
	p.Mode = internal.WhisperMode(r.ReadU8())
	p.Target = r.ReadStr16()
	if p.Mode == internal.WhisperMode_WHISPER_MODE_SEND {
		p.Message = r.ReadStr16()
	}
}
