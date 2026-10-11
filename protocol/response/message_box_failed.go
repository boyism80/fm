package response

import "github.com/boyism80/fm/stream"

type MessageBoxFailed struct{}

func (*MessageBoxFailed) Opcode() uint16 { return 0xC8 }

func (*MessageBoxFailed) Serialize(*stream.StreamWriter) error { return nil }

func (*MessageBoxFailed) Deserialize(*stream.StreamReader) {}
