package conn

import (
	"fmt"

	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type packet interface {
	Opcode() uint16
	Deserialize(reader *stream.StreamReader)
}

type decoder struct {
	new  func() packet
	peek func(body []byte) bool
}

func newDecoder[T any, P interface {
	*T
	packet
}](peek func(body []byte) bool) decoder {
	return decoder{new: func() packet { return P(new(T)) }, peek: peek}
}

var decoders = func() []decoder {
	dialog := func(types ...constant.DialogType) func([]byte) bool {
		return func(body []byte) bool {
			if len(body) < 6 {
				return false
			}
			for _, t := range types {
				if constant.DialogType(body[5]) == t {
					return true
				}
			}
			return false
		}
	}
	party := func(sub pconst.PartySubOpcode) func([]byte) bool {
		return func(body []byte) bool {
			return len(body) > 0 && pconst.PartySubOpcode(body[0]) == sub
		}
	}
	partyUpdate := func(match func(body []byte) bool) func([]byte) bool {
		return func(body []byte) bool {
			return len(body) >= 11 && pconst.PartySubOpcode(body[0]) == pconst.PartyS2CPartyUpdate && match(body)
		}
	}

	return []decoder{
		newDecoder[response.Authenticate](func(body []byte) bool { return len(body) > 0 && body[0] == 0 }),
		newDecoder[response.LoginFailed](func(body []byte) bool { return len(body) > 0 }),
		newDecoder[response.EndOfServerList](func(body []byte) bool { return len(body) > 0 && body[0] == 0xFF }),
		newDecoder[response.ServerList](nil),
		newDecoder[response.CharacterList](nil),
		newDecoder[response.Transfer](nil),
		newDecoder[response.CheckName](nil),
		newDecoder[response.CreateCharacter](nil),
		newDecoder[response.Ping](nil),
		newDecoder[response.Warp](func(body []byte) bool { return len(body) > 4 && body[4] == 2 }),
		newDecoder[response.FieldRelocate](nil),
		newDecoder[response.Login](nil),
		newDecoder[response.KeyMap](nil),
		newDecoder[response.Notice](nil),
		newDecoder[response.SpawnNpc](nil),
		newDecoder[response.RemoveNpc](nil),
		newDecoder[response.DialogYesNo](dialog(constant.DialogTypeYesNo)),
		newDecoder[response.DialogInput](dialog(constant.DialogTypeInput)),
		newDecoder[response.DialogList](dialog(constant.DialogTypeList)),
		newDecoder[response.DialogStyle](dialog(constant.DialogTypeStyle)),
		newDecoder[response.DialogAccept](dialog(constant.DialogTypeAccept, constant.DialogTypeAcceptEscape)),
		newDecoder[response.Dialog](func(body []byte) bool { return len(body) >= 6 }),
		newDecoder[response.UpdateStats](nil),
		newDecoder[response.OpenNpcShop](nil),
		newDecoder[response.InventoryOperation](nil),
		newDecoder[response.UpdateQuest](func(body []byte) bool { return len(body) > 0 && body[0] == 1 }),
		newDecoder[response.GuildInvite](func(body []byte) bool {
			return len(body) > 0 && pconst.GuildSubOpcode(body[0]) == pconst.GuildS2CInvite
		}),
		newDecoder[response.GuildMessage](nil),
		newDecoder[response.SwitchChannel](nil),
		newDecoder[response.ServerBlocked](nil),
		newDecoder[response.PartyCreated](party(pconst.PartyS2CPartyCreated)),
		newDecoder[response.PartyInvite](party(pconst.PartyS2CInvite)),
		newDecoder[response.PartyUpdateJoin](party(pconst.PartyS2CPartyJoin)),
		newDecoder[response.PartyUpdateLeaderChange](party(pconst.PartyS2CLeaderChange)),
		newDecoder[response.PartyUpdateDisband](partyUpdate(func(body []byte) bool { return body[9] == 0 })),
		newDecoder[response.PartyUpdateExpel](partyUpdate(func(body []byte) bool { return body[10] == 1 })),
		newDecoder[response.PartyUpdateLeave](partyUpdate(func([]byte) bool { return true })),
		newDecoder[response.PartyStatusMessage](func(body []byte) bool {
			if len(body) == 0 {
				return false
			}
			switch pconst.PartySubOpcode(body[0]) {
			case pconst.PartyS2CPartyUpdate, pconst.PartyS2CSilentUpdate, pconst.PartyS2CPartyPortal:
				return false
			}
			return true
		}),
		newDecoder[response.SpawnMob](nil),
		newDecoder[response.DieMob](nil),
		newDecoder[response.SpawnMeso](func(body []byte) bool { return len(body) >= 6 && body[5] == 1 }),
		newDecoder[response.SpawnItem](func(body []byte) bool { return len(body) >= 6 }),
		newDecoder[response.RemoveItem](nil),
		newDecoder[response.ShowBossHp](func(body []byte) bool {
			return len(body) > 0 && response.EnvironmentChangeMode(body[0]) == response.EnvironmentChangeModeBossHP
		}),
		newDecoder[response.EnvironmentChange](func(body []byte) bool { return len(body) > 0 }),
		newDecoder[response.SpawnReactor](nil),
		newDecoder[response.TriggerReactor](nil),
		newDecoder[response.DestroyReactor](nil),
		newDecoder[response.CarnivalStart](nil),
		newDecoder[response.CarnivalObtainedCP](nil),
		newDecoder[response.CarnivalPartyCP](nil),
		newDecoder[response.CarnivalSummon](nil),
		newDecoder[response.CarnivalDied](nil),
	}
}()

var decodersByOpcode = func() map[uint16][]decoder {
	m := make(map[uint16][]decoder)
	for _, d := range decoders {
		op := d.new().Opcode()
		m[op] = append(m[op], d)
	}
	return m
}()

var Responses = func() []any {
	out := make([]any, 0, len(decoders))
	for _, d := range decoders {
		out = append(out, d.new())
	}
	return out
}()

func Decode(opcode uint16, body []byte) (any, error) {
	for _, d := range decodersByOpcode[opcode] {
		if d.peek != nil && d.peek(body) == false {
			continue
		}
		pkt := d.new()
		pkt.Deserialize(stream.NewStreamReader(&body, stream.LittleEndian))
		return pkt, nil
	}
	return nil, fmt.Errorf("opcode %d is not decoded", opcode)
}
