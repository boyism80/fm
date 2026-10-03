package response

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/boyism80/fm/protocol/dto"
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

func TestWelcomeBytes(t *testing.T) {
	p := &Welcome{RecvIv: []byte{0x65, 0x56, 0x12, 0xFD}, SendIv: []byte{0x2F, 0xA3, 0x65, 0x43}}
	w := stream.NewStreamWriter(stream.LittleEndian)
	if err := p.Serialize(w); err != nil {
		t.Fatal(err)
	}
	body := w.Bytes()
	if binary.LittleEndian.Uint16(body[:2]) != MAGIC {
		t.Fatalf("magic=%x", body[:2])
	}
	if p.Opcode() != uint16(len(body)) {
		t.Fatalf("length header=%d body=%d", p.Opcode(), len(body))
	}
	if body[len(body)-1] != 1 {
		t.Fatalf("tail=%x", body[len(body)-1])
	}
	roundTrip(t, p)
}

func TestLoginFailedBytes(t *testing.T) {
	p := &LoginFailed{Reason: LoginFailedReasonNoPopup}
	w := stream.NewStreamWriter(stream.LittleEndian)
	if err := p.Serialize(w); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(w.Bytes(), []byte{LoginFailedReasonNoPopup}) {
		t.Fatalf("body=%x", w.Bytes())
	}
	roundTrip(t, p)
	roundTrip(t, &LoginFailed{Reason: LoginFailedReasonAlreadyLoggedIn})
	roundTrip(t, &LoginFailed{Reason: LoginFailedReasonPasswordChangeRequired})
}

func TestLoginResponseRoundTrip(t *testing.T) {
	roundTrip(t, &Authenticate{
		AccountId: 9, Gender: 1, Role: 1, AccountName: "bot",
		IsChatBlocked: true, ChatBlockTime: 99,
	})
	roundTrip(t, &ServerList{
		ServerId: 0, WorldName: "Scania", Flag: 1, EventMessage: "hello",
		Channels: []ServerChannel{{ChannelID: 0, Name: "Ch.1", Load: 10}},
	})
	roundTrip(t, &EndOfServerList{})
	roundTrip(t, &CharacterList{
		SlotCount: 3,
		Characters: []dto.Character{{
			ID: 7, Name: "홍길동", Gender: 1, SkinColor: 2,
			Face: 20000, Hair: 30000, Level: 10, Class: 100,
			Str: 4, Dex: 4, Int: 4, Luk: 4,
			Hp: 50, MaxHp: 50, Mp: 20, MaxMp: 20,
			AbilityPoint: 1, Exp: 100, Population: 1,
			Map: 100000000, SpawnPoint: 2, Mega: true,
			BaseLooks: map[int8]uint32{1: 1040002},
			Overlays:  map[int8]uint32{5: 1060002},
			Weapon:    1302000,
			Rank:      3, RankDiff: -1, ClassRank: 4, ClassRankDiff: 2,
		}},
	})
	roundTrip(t, &CheckName{Name: "Bot", Exists: true})
	roundTrip(t, &CreateCharacter{
		Success: true,
		Character: &dto.Character{
			ID: 1, Name: "Bot", Gender: 0, Face: 20000, Hair: 30000, Level: 1,
		},
	})
	roundTrip(t, &Transfer{IP: "127.0.0.1", Port: 8485, CharacterId: 7})
	roundTrip(t, &Ping{})
}
