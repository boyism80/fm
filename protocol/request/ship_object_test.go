package request

import (
	"bytes"
	"testing"

	"github.com/boyism80/fm/stream"
)

func TestShipObjectRoundTrip(t *testing.T) {
	p := &ShipObject{MapID: 101000300}
	if p.Opcode() != 0xB3 {
		t.Fatalf("opcode=%x want=0xB3", p.Opcode())
	}
	w := stream.NewStreamWriter(stream.LittleEndian)
	if err := p.Serialize(w); err != nil {
		t.Fatal(err)
	}
	want := []byte{0x6C, 0x24, 0x05, 0x06}
	if !bytes.Equal(w.Bytes(), want) {
		t.Fatalf("body=%x want=%x", w.Bytes(), want)
	}
	data := w.Bytes()
	var out ShipObject
	out.Deserialize(stream.NewStreamReader(&data, stream.LittleEndian))
	if out.MapID != 101000300 {
		t.Fatalf("map=%d", out.MapID)
	}
}
