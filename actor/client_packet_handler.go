package actor

import (
	"log"

	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/packet/req"
	"github.com/boyism80/fm/packet/resp"
	"github.com/boyism80/fm/types"
)

func RegisterPacketHandler(client *ClientActor, h *handler.PacketHandler) {
	handler.RegisterPacketHandler(0x0A, client, h, onPong)
	handler.RegisterPacketHandler(0x01, client, h, onLogin)
	handler.RegisterPacketHandler(0x04, client, h, onCharacterList)
	handler.RegisterPacketHandler(0x07, client, h, onCheckName)
	handler.RegisterPacketHandler(0x08, client, h, onCreateCharacter)
	handler.RegisterPacketHandler(0x09, client, h, onDeleteCharacter)
	handler.RegisterPacketHandler(0x05, client, h, onSelectCharacter)

	// game
	handler.RegisterPacketHandler(0x06, client, h, onLoginGame)
}

func onPong(client *ClientActor, request *req.Pong) {
}

func onLogin(client *ClientActor, request *req.Login) {

	if request.Id == "cshyeon" {
		client.Send(&resp.Authenticate{
			AccountId:     2390,
			Gender:        0,
			Admin:         true,
			AccountName:   request.Id,
			IsChatBlocked: false,
			ChatBlockTime: 116445060000000000,
		}, types.SEND_POLICY_ENCRYPT)

		for i := 0; i < 10; i++ {
			client.Send(&resp.ServerList{
				ServerId:     uint8(i),
				ChannelSize:  5,
				WorldName:    "cshyeon",
				Flag:         0,
				EventMessage: "채승현이 간다.",
			}, types.SEND_POLICY_ENCRYPT)
		}

		client.Send(&resp.EndOfServerList{}, types.SEND_POLICY_ENCRYPT)
	} else {
		client.Send(&resp.LoginFailed{Reason: resp.LoginFailedReasonNoPopup}, types.SEND_POLICY_ENCRYPT)
		client.Send(&resp.Notice{
			Type:    resp.NoticeTypePopup,
			Channel: 0,
			Message: "Hello",
			MegaEar: false}, types.SEND_POLICY_ENCRYPT)
	}
}

func onCharacterList(client *ClientActor, request *req.CharacterList) {
	client.Send(&resp.CharacterList{
		Characters: []entity.Character{
			entity.NewDummyCharacter(1, "채승현"),
			entity.NewDummyCharacter(2, "채진영"),
		},
		SlotCount: 6,
	}, types.SEND_POLICY_ENCRYPT)
}

func onCheckName(client *ClientActor, request *req.CheckName) {
	exists := request.Name == "채승현"
	client.Send(&resp.CheckName{Name: request.Name, Exists: exists}, types.SEND_POLICY_ENCRYPT)
}

func onCreateCharacter(client *ClientActor, request *req.CreateCharacter) {

	success := request.Name != "채진영"
	client.Send(&resp.CreateCharacter{
		Success: success,
		Character: &entity.Character{
			Id:         1,
			Name:       request.Name,
			Gender:     0,
			SkinColor:  0,
			Face:       request.Face,
			Hair:       request.Hair,
			Level:      1,
			Class:      0,
			Str:        12,
			Dex:        5,
			Int:        4,
			Luk:        4,
			Hp:         50,
			MaxHp:      50,
			Mp:         5,
			MaxMp:      5,
			SpawnPoint: 3,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func onDeleteCharacter(client *ClientActor, request *req.DeleteCharacter) {
	client.Send(&resp.DeleteCharacter{
		Id:      request.Id,
		Success: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func onSelectCharacter(client *ClientActor, request *req.SelectCharacter) {
	client.Send(&resp.Transfer{
		IP:          "127.0.0.1",
		Port:        7100,
		CharacterId: request.CharacterId,
	}, types.SEND_POLICY_ENCRYPT)
}

func onLoginGame(client *ClientActor, request *req.LoginGame) {
	ch := entity.NewDummyCharacter(35177, "강원기")
	m := entity.NewMap()
	m.MapObjects[ch.Id] = &ch

	from := ch.Position()
	var pk types.Packet
	for _, obj := range m.GetMapObjectsInRange(from, types.MaxViewRangeSq, types.ObjectTypeObject) {
		switch obj.Type() {
		case types.ObjectTypePlayer:
			character, ok := obj.(*entity.Character)
			if !ok {
				log.Printf("Object type mismatch: expected *Character, got %T", obj)
				continue
			}
			pk = &resp.SpawnPlayer{
				Character:       character,
				BuffStates:      [4]uint32{0, 0, 0, 0},
				Diseases:        [4]uint32{0, 0, 0, 0},
				Position:        character.Position(),
				CrushRings:      []*entity.Ring{},
				FriendshipRings: []*entity.Ring{},
				MarriageRings:   []*entity.Ring{},
			}

		default:
			pk = nil
		}

		if pk != nil {
			ch.Send(pk, types.SEND_POLICY_ENCRYPT)
		}
	}

	client.Send(&resp.Warp{Character: &ch}, types.SEND_POLICY_ENCRYPT)
}
