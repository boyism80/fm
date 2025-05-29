package actor

import (
	"log"
	"math"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/crypt"
	"github.com/boyism80/fm/common/handler"
	common_msg "github.com/boyism80/fm/common/msg"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/common/util"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol/resp"

	common_resp "github.com/boyism80/fm/common/protocol/resp"
	login_resp "github.com/boyism80/fm/login/protocol/resp"
)

func RegisterGameClientMessageHandlers(ctx protoactor.Context, m *GameClientActor, h *handler.MessageHandler) {
	handler.RegisterHandler(ctx, m, h, onGameClientConnected)
	handler.RegisterHandler(ctx, m, h, onGameClientStopping)
	handler.RegisterHandler(ctx, m, h, onGameClientInvokeHandler)
	handler.RegisterHandler(ctx, m, h, onGameClientSendProtocol)
	handler.RegisterHandler(ctx, m, h, onGameClientPing)
	handler.RegisterHandler(ctx, m, h, onGameClientWarped)
	handler.RegisterHandler(ctx, m, h, onGameClientItemLooting)
	handler.RegisterHandler(ctx, m, h, onGameClientMesoLooting)
}

func (client *GameClientActor) ReceivePackets(ctx protoactor.Context, pid *protoactor.PID) {
	buf := make([]byte, 1024)
	conn := client.conn

	for {
		n, err := conn.Read(buf)
		if err != nil {
			ctx.Stop(ctx.Self())
			break
		}

		client.buffer = append(client.buffer, buf[:n]...)
		reader := stream.NewStreamReader(&client.buffer, stream.LittleEndian)

		for {
			header, err := reader.Read(4)
			if err != nil {
				break
			}

			if !client.receiveCrypt.CheckPacketHeader(header) {
				break
			}

			len, err := crypt.GetPacketLength(header)
			if err != nil {
				break
			}

			data, err := reader.Read(int(len))
			if err != nil {
				break
			}

			data = client.receiveCrypt.Decrypt(data)
			log.Println("[R] " + util.ToHexString(data))
			opcodeReader := stream.NewStreamReader(&data, stream.LittleEndian)
			opcode, err := opcodeReader.ReadU16()
			if err != nil {
				break
			}

			data = data[2:]
			ctx.Send(pid, &common_msg.InvokeHandler{
				Opcode: opcode,
				Data:   data,
			})
			reader.DiscardRead()
		}
	}
}

func onGameClientConnected(ctx protoactor.Context, client *GameClientActor, m *common_msg.ClientConnected) {
	log.Println("클라이언트 연결됨")
	// timer := scheduler.NewTimerScheduler(ctx.ActorSystem().Root)
	// client.stopTimer = timer.SendRepeatedly(10*time.Second, 10*time.Second, ctx.Self(), &msg.Ping{})

	// 게임서버에서 보내주면 안됨
	client.Send(&common_resp.Welcome{
		SendIv: client.sendCrypt.IV(),
		RecvIv: client.receiveCrypt.IV()}, types.SEND_POLICY_RAW)

	client.Send(&login_resp.LoginFailed{
		Reason: login_resp.LoginFailedReasonNoPopup}, types.SEND_POLICY_ENCRYPT)
	go client.ReceivePackets(ctx, ctx.Self())
}

func onGameClientStopping(ctx protoactor.Context, client *GameClientActor, m *protoactor.Stopping) {
	ch := client.ch
	if ch != nil {
		mapActor := client.ctx.MapActors[client.ch.Map]
		if mapActor != nil {
			ctx.Send(mapActor, &msg.LeaveMap{
				ID: ch.ID,
			})
		}
	}

	log.Println("클라이언트 접속 종료")
	client.conn.Close()
	if client.stopTimer != nil {
		client.stopTimer()
	}
}

func onGameClientInvokeHandler(ctx protoactor.Context, client *GameClientActor, m *common_msg.InvokeHandler) {
	client.Invoke(ctx, int(m.Opcode), m.Data)
}

func onGameClientSendProtocol(ctx protoactor.Context, state *GameClientActor, m *common_msg.SendProtocol) {
	state.Send(m.Protocol, m.Policy)
}

func onGameClientPing(ctx protoactor.Context, client *GameClientActor, m *common_msg.Ping) {
	client.Send(&common_resp.Ping{}, types.SEND_POLICY_ENCRYPT)
}

func onGameClientWarped(ctx protoactor.Context, client *GameClientActor, m *msg.Warped) {
	// 어떤 플레이어가 내가 속한 맵으로 왔다.
	// 나의 캐릭터 정보를 해당 플레이어에게 보내준다.

	ctx.Send(m.Sender, &common_msg.SendProtocol{
		Protocol: &resp.SpawnPlayer{
			Character:       client.ch,
			BuffStates:      [4]uint32{},
			Diseases:        [4]uint32{},
			CrushRings:      []*entity.Ring{},
			FriendshipRings: []*entity.Ring{},
			MarriageRings:   []*entity.Ring{},
		},
		Policy: types.SEND_POLICY_ENCRYPT,
	})
}

func onGameClientItemLooting(ctx protoactor.Context, client *GameClientActor, m *msg.CharacterItemLooting) {

	invenType := m.Item.GetInventoryType()
	inven := client.ch.Inventory[invenType]
	spec := m.Item.GetSpec()
	gain := uint16(0)
	if !inven.IsFree(spec, m.Item.GetCount()) {
		client.Send(&resp.ItemGainFailed{
			Mode: resp.ItemGainFailedTypeFull,
		}, types.SEND_POLICY_ENCRYPT)

		client.Send(&resp.UpdateStats{
			UnlockAction: true,
		}, types.SEND_POLICY_ENCRYPT)

		ctx.Send(m.PID, &msg.ItemLooted{
			Success:     false,
			Count:       0,
			CharacterId: client.ch.ID,
		})
		return
	}

	for m.Item.GetCount() > 0 {
		slot, ok := inven.FindSlot(spec)
		if !ok {
			break
		}

		exists, ok := inven.Items[int16(slot)]
		cap := uint16(0)
		if ok {
			// 기존에 있는 아이템 수량 증가
			cap = min(spec.GetCapacity()-exists.GetCount(), m.Item.GetCount())
			exists.Increase(cap)
			client.Send(&resp.UpdateInventorySlot{
				InventoryType: invenType,
				Slot:          int16(slot),
				Item:          exists,
			}, types.SEND_POLICY_ENCRYPT)
		} else {
			// 새로운 슬롯에 아이템 추가
			cap = min(spec.GetCapacity(), m.Item.GetCount())
			inven.Items[int16(slot)] = m.Item.Clone(cap)
			client.Send(&resp.AddInventorySlot{
				InventoryType: invenType,
				Slot:          int16(slot),
				Item:          inven.Items[int16(slot)],
			}, types.SEND_POLICY_ENCRYPT)
		}
		m.Item.Reduce(cap)
		gain += cap
	}

	client.Send(&resp.ShowItemGain{
		ItemId: spec.GetID(),
		Count:  uint32(gain),
		Mode:   resp.ShowItemGainTypeStatus,
	}, types.SEND_POLICY_ENCRYPT)

	ctx.Send(m.PID, &msg.ItemLooted{
		Success:     true,
		Count:       int32(gain),
		CharacterId: client.ch.ID,
	})
}

func onGameClientMesoLooting(ctx protoactor.Context, client *GameClientActor, m *msg.CharacterMesoLooting) {
	cap := math.MaxInt32 - client.ch.Meso
	if m.Meso > cap {
		client.Send(&resp.ItemGainFailed{
			Mode: resp.ItemGainFailedTypeFull,
		}, types.SEND_POLICY_ENCRYPT)

		client.Send(&resp.UpdateStats{
			UnlockAction: true,
		}, types.SEND_POLICY_ENCRYPT)

		ctx.Send(m.PID, &msg.ItemLooted{
			Success:     false,
			Count:       0,
			CharacterId: client.ch.ID,
		})
		return
	}

	client.ch.Meso += m.Meso
	client.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.StatMeso: client.ch.Meso,
		},
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)

	client.Send(&resp.ShowMesoGain{
		Count: m.Meso,
		Mode:  resp.ShowMesoGainTypeStatus,
	}, types.SEND_POLICY_ENCRYPT)

	ctx.Send(m.PID, &msg.ItemLooted{
		Success:     true,
		Count:       m.Meso,
		CharacterId: client.ch.ID,
	})
}
