package request

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/boyism80/fm/stream"
)

func roundTrip(t *testing.T, pkt interface {
	Serialize(*stream.StreamWriter) error
	Deserialize(*stream.StreamReader)
}) {
	t.Helper()
	w := stream.NewStreamWriter(stream.LittleEndian)
	if err := pkt.Serialize(w); err != nil {
		t.Fatal(err)
	}
	data := w.Bytes()
	out := reflect.New(reflect.TypeOf(pkt).Elem()).Interface().(interface {
		Deserialize(*stream.StreamReader)
	})
	out.Deserialize(stream.NewStreamReader(&data, stream.LittleEndian))
	if !reflect.DeepEqual(pkt, out) {
		t.Fatalf("got %#v want %#v", out, pkt)
	}
}

func TestLoginRoundTrip(t *testing.T) {
	p := &Login{ID: "ab", Pw: "pw", Mac: "AA-BB-CC-DD-EE-FF"}
	w := stream.NewStreamWriter(stream.LittleEndian)
	if err := p.Serialize(w); err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x02, 0x00, 'a', 'b',
		0x02, 0x00, 'p', 'w',
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
	}
	if !bytes.Equal(w.Bytes(), want) {
		t.Fatalf("body=%x want=%x", w.Bytes(), want)
	}
	roundTrip(t, p)
}

func TestLoginCodecRoundTrip(t *testing.T) {
	roundTrip(t, &CharacterList{Server: 1, Channel: 0})
	roundTrip(t, &CheckName{Name: "홍길동"})
	roundTrip(t, &CreateCharacter{
		Name: "Bot", Face: 20000, Hair: 30000,
		Top: 1040002, Bottom: 1060002, Shoes: 1072001, Weapon: 1302000,
	})
	roundTrip(t, &SelectCharacter{CharacterId: 42})
	roundTrip(t, &DeleteCharacter{ID: 42})
	roundTrip(t, &LoginGame{PlayerId: 42})
	roundTrip(t, &Pong{})
}
