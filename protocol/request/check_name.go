package request

import "github.com/boyism80/fm/stream"

type CheckName struct {
	Name string
}

func (a *CheckName) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *CheckName) Deserialize(reader *stream.StreamReader) {
	a.Name = reader.ReadStr16()
}
