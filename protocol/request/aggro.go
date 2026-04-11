package request

import (
	"github.com/boyism80/fm/stream"
)

type Aggro struct{}

func (p *Aggro) Serialize(writer *stream.StreamWriter) error {
	// TODO: implement after packet structure is analyzed
	return nil
}

func (p *Aggro) Deserialize(reader *stream.StreamReader) {
	// TODO: implement after packet structure is analyzed
}
