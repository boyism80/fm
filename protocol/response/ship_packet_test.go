package response

import (
	"bytes"
	"testing"

	"github.com/boyism80/fm/stream"
)

func TestShipSpecialEffectSerialize(t *testing.T) {
	p := &ShipSpecialEffect{Effect: ShipSpecialBalrog}
	w := stream.NewStreamWriter(stream.LittleEndian)
	if err := p.Serialize(w); err != nil {
		t.Fatal(err)
	}
	want := []byte{0x0A, 0x04}
	if !bytes.Equal(w.Bytes(), want) {
		t.Fatalf("body=%x want=%x", w.Bytes(), want)
	}
	if p.Opcode() != 0x66 {
		t.Fatalf("opcode=%x want=0x66", p.Opcode())
	}
}

func TestShipStateSerialize(t *testing.T) {
	cases := []struct {
		name  string
		state uint16
		want  []byte
	}{
		{"leaving", ShipStateLeaving, []byte{0x03, 0x00}},
		{"docked", ShipStateDocked, []byte{0x01, 0x00}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &ShipState{State: tc.state}
			w := stream.NewStreamWriter(stream.LittleEndian)
			if err := p.Serialize(w); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(w.Bytes(), tc.want) {
				t.Fatalf("body=%x want=%x", w.Bytes(), tc.want)
			}
			if p.Opcode() != 0x67 {
				t.Fatalf("opcode=%x want=0x67", p.Opcode())
			}
		})
	}
}
