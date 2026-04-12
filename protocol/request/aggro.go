package request

import (
	"github.com/boyism80/fm/stream"
)

type Aggro struct{}

func (p *Aggro) Serialize(writer *stream.StreamWriter) error {

	return nil
}

func (p *Aggro) Deserialize(reader *stream.StreamReader) {

}
