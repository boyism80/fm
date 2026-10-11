package request

import "github.com/boyism80/fm/stream"

type CloseChalkboard struct{}

func (*CloseChalkboard) Opcode() byte { return 0x21 }

func (*CloseChalkboard) Serialize(*stream.StreamWriter) error { return nil }

func (*CloseChalkboard) Deserialize(*stream.StreamReader) {}
