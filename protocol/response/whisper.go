package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type Whisper struct {
	Mode     internal.WhisperMode
	Name     string
	Message  string
	Success  bool
	Location pconst.WhisperLocation
	Channel  int32
	MapID    uint32
}

func (p *Whisper) Opcode() uint16 {
	return 0x5D
}

func (p *Whisper) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Mode))
	w.WriteStr16(p.Name)
	switch p.Mode {
	case internal.WhisperMode_WHISPER_MODE_SEND_RESULT:
		w.WriteBoolean(p.Success)
	case internal.WhisperMode_WHISPER_MODE_RECEIVE:
		w.Write16(int16(p.Channel))
		w.WriteStr16(p.Message)
	case internal.WhisperMode_WHISPER_MODE_FIND_RESULT, internal.WhisperMode_WHISPER_MODE_FIND_BUDDY_RESULT:
		w.WriteU8(uint8(p.Location))
		if p.Location == pconst.WhisperLocationMap {
			w.WriteU32(p.MapID)
			w.WriteU64(0)
		} else {
			w.Write32(p.Channel)
		}
	}
	return nil
}

func (p *Whisper) Deserialize(r *stream.StreamReader) {
	p.Mode = internal.WhisperMode(r.ReadU8())
	p.Name = r.ReadStr16()
	switch p.Mode {
	case internal.WhisperMode_WHISPER_MODE_SEND_RESULT:
		p.Success = r.ReadBool()
	case internal.WhisperMode_WHISPER_MODE_RECEIVE:
		p.Channel = int32(r.Read16())
		p.Message = r.ReadStr16()
	case internal.WhisperMode_WHISPER_MODE_FIND_RESULT, internal.WhisperMode_WHISPER_MODE_FIND_BUDDY_RESULT:
		p.Location = pconst.WhisperLocation(r.ReadU8())
		if p.Location == pconst.WhisperLocationMap {
			p.MapID = r.ReadU32()
			r.Skip(8)
		} else {
			p.Channel = r.Read32()
		}
	}
}
