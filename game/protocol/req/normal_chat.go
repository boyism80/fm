package req

import (
	"github.com/boyism80/fm/common/stream"
)

type NormalChat struct {
	Message string
	Show    bool
}

func (m *NormalChat) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *NormalChat) Deserialize(reader *stream.StreamReader) error {
	message, err := reader.ReadStr16()
	if err != nil {
		return err
	}

	show, err := reader.ReadBool()
	if err != nil {
		return err
	}

	m.Message = message
	m.Show = show
	return nil
}
