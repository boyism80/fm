package request

import "github.com/boyism80/fm/stream"

type CheckName struct {
	Name string
}

func (a *CheckName) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *CheckName) Deserialize(reader *stream.StreamReader) error {
	name, err := reader.ReadStr16()
	if err != nil {
		return err
	}
	a.Name = name
	return nil
}
