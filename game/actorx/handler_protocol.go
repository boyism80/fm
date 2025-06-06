package actorx

import (
	"log"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/common/luax"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol"
	"github.com/boyism80/fm/game/protocol/req"
	"github.com/boyism80/fm/game/protocol/resp"
	lua "github.com/yuin/gopher-lua"

	common_msg "github.com/boyism80/fm/common/msg"
	common_req "github.com/boyism80/fm/common/protocol/req"
)

func RegisterGameClientPacketHandler(ctx actor.Context, client *GameClientActor, h *handler.PacketHandler) {
	handler.RegisterPacketHandler(0x0A, ctx, client, h, onGameClientPong)
	handler.RegisterPacketHandler(0x06, ctx, client, h, onLoginGame)
	handler.RegisterPacketHandler(0x18, ctx, client, h, onGameClientMovePlayer)
	handler.RegisterPacketHandler(0x20, ctx, client, h, onGameClientNormalChat)
	handler.RegisterPacketHandler(0x1B, ctx, client, h, onGameClientAttack)
	handler.RegisterPacketHandler(0x36, ctx, client, h, onGameMoveItem)
	handler.RegisterPacketHandler(0x34, ctx, client, h, onGameSortInventory)
	handler.RegisterPacketHandler(0xA3, ctx, client, h, onGameItemLoot)
	handler.RegisterPacketHandler(0x4D, ctx, client, h, onGameDropMeso)
	handler.RegisterPacketHandler(0x15, ctx, client, h, onGameWarp)
	handler.RegisterPacketHandler(0x9E, ctx, client, h, onGameClientNpcControl)
	handler.RegisterPacketHandler(0x2B, ctx, client, h, onGameClientDialog)
	handler.RegisterPacketHandler(0x29, ctx, client, h, onGameClientNpcClick)
}

func onGameClientPong(ctx actor.Context, client *GameClientActor, request *common_req.Pong) {
}

func onLoginGame(ctx actor.Context, client *GameClientActor, request *req.LoginGame) {
	name := "채승현"
	if request.PlayerId != 1 {
		name = "채진영"
	}
	ch := entity.NewDummyCharacter(client, client.listener, request.PlayerId, name, client.serverContext)
	client.ch = &ch
	RegisterLifeHandlers(ctx, &ch.Life, client.messageHandler)

	mapActor := client.serverContext.MapActors[client.ch.Map]
	if mapActor != nil {
		ctx.Send(mapActor, &msg.EnterMap{
			ID:         ch.ID,
			PID:        ctx.Self(),
			Init:       true,
			SpawnPoint: ch.SpawnPoint,
		})
	}
}

func onGameClientMovePlayer(ctx actor.Context, client *GameClientActor, req *req.MovePlayer) {

	beforePosition := client.ch.Position

	for _, frag := range req.Fragments {
		if move, ok := frag.(*protocol.AbsoluteLifeMovement); ok {
			client.ch.Position = move.Position
		}

		client.ch.Stance = frag.GetStance()
	}

	mapActor := client.serverContext.MapActors[client.ch.Map]
	if mapActor == nil {
		return
	}

	ctx.Send(mapActor, &msg.MapBroadcastRange{
		Sender: ctx.Self(),
		Pivot:  beforePosition,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.Move{
				Character:  client.ch,
				Fragments:  req.Fragments,
				StartPoint: beforePosition,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}

func onGameClientNormalChat(ctx actor.Context, client *GameClientActor, req *req.NormalChat) {

	if strings.HasPrefix(req.Message, "/") {
		params := strings.Split(strings.TrimPrefix(req.Message, "/"), " ")
		err := client.commandHandler.Handle(ctx, params...)
		if err == nil {
			return
		}
	}

	client.listener.OnChat(req.Message, false, req.DontRecordHistory)
}

func onGameClientAttack(ctx actor.Context, client *GameClientActor, req *req.Attack) {
	position := client.ch.Position
	mapActor := client.serverContext.MapActors[client.ch.Map]
	if mapActor == nil {
		return
	}

	ctx.Send(mapActor, &msg.MapBroadcastRange{
		Sender: ctx.Self(),
		Pivot:  position,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.Attack{
				CharacterId: client.ch.ID,
				AttackInfo:  req.AttackInfo,
				SkillLevel:  0,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}

func onGameMoveItem(ctx actor.Context, client *GameClientActor, req *req.MoveItem) {
	if req.Source < 0 {
		client.Unequip(ctx, constant.EquipmentPartsType(req.Source), req.Dest)
	} else if req.Dest < 0 {
		client.Equip(ctx, constant.EquipmentPartsType(req.Dest), req.Source)
	} else if req.Dest == 0 {
		client.Drop(ctx, req.InventoryType, req.Source, req.Count)
	} else {
		client.MoveItem(req.InventoryType, req.Source, req.Dest)
	}
}

func onGameSortInventory(ctx actor.Context, client *GameClientActor, req *req.SortInventory) {
	client.MergeItems(req.InventoryType)
	client.SortInventory(req.InventoryType)

	client.Send(&resp.EndSortInventory{
		InventoryType: req.InventoryType,
	}, types.SEND_POLICY_ENCRYPT)

	client.Send(&resp.UpdateStats{
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func onGameItemLoot(ctx actor.Context, client *GameClientActor, req *req.ItemLoot) {

	mpid, ok := client.serverContext.MapActors[client.ch.Map]
	if !ok {
		client.Send(&resp.UpdateStats{
			UnlockAction: true,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}

	ctx.Send(mpid, &msg.MapItemLoot{
		Actor:       ctx.Self(),
		OID:         req.OID,
		CharacterId: client.ch.ID,
		Position:    req.Position,
	})
}

func onGameDropMeso(ctx actor.Context, client *GameClientActor, req *req.DropMeso) {
	ch := client.ch
	mapActor := client.serverContext.MapActors[ch.Map]
	if mapActor == nil {
		return
	}

	if req.Count < 10 || req.Count > 50000 {
		ctx.Stop(ctx.Self())
		return
	}

	if req.Count > ch.Meso {
		client.Send(&resp.UpdateStats{
			UnlockAction: true,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}

	ch.Meso -= req.Count
	client.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.StatMeso: ch.Meso,
		},
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)

	ctx.Send(mapActor, &msg.MapSpawnMeso{
		Count:        req.Count,
		SpawnedPoint: ch.Position,
		Owner:        ctx.Self(),
		OwnerID:      ch.ID,
	})
}

func onGameWarp(ctx actor.Context, client *GameClientActor, req *req.Warp) {
	if req.Target != 0xFFFFFFFF {
		return
	}
	oldMapSpec, ok := client.serverContext.Resources.Maps[client.ch.Map]
	if !ok {
		return
	}

	oldPortal, ok := oldMapSpec.FindPortal(req.PortalName)
	if !ok {
		client.Send(&resp.UpdateStats{
			UnlockAction: true,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}

	newMapSpec, ok := client.serverContext.Resources.Maps[uint32(oldPortal.TargetMapId)]
	if !ok {
		client.Send(&resp.UpdateStats{
			UnlockAction: true,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}

	newPortal, ok := newMapSpec.FindPortal(oldPortal.Target)
	if !ok {
		client.Send(&resp.UpdateStats{
			UnlockAction: true,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}

	client.Warp(ctx, uint32(oldPortal.TargetMapId), newPortal.ID)
}

func onGameClientNpcControl(ctx actor.Context, client *GameClientActor, req *req.NpcAction) {

	client.Send(&resp.NpcAction{
		Bytes: req.Bytes,
	}, types.SEND_POLICY_ENCRYPT)
}

func onGameClientDialog(ctx actor.Context, client *GameClientActor, req *req.Dialog) {

	if client.ch.Dialog == nil {
		return
	}

	args := []lua.LValue{}

	co := client.ch.Dialog
	client.ch.Dialog = nil

	switch req.DialogType {
	case constant.DIALOG_TYPE_DEFAULT:
		if req.Next {
			args = append(args, lua.LTrue)
		} else {
			args = append(args, lua.LFalse)
		}

	case constant.DIALOG_TYPE_YES_NO:
		if req.Next {
			args = append(args, lua.LTrue)
		} else {
			args = append(args, lua.LFalse)
		}

	case constant.DIALOG_TYPE_LIST:
		if req.Next {
			args = append(args, lua.LNumber(req.Selected))
		} else {
			args = append(args, lua.LNil)
		}

	case constant.DIALOG_TYPE_INPUT:
		if req.Next {
			args = append(args, lua.LString(req.Text))
		} else {
			args = append(args, lua.LNil)
		}

	case constant.DIALOG_TYPE_ACCEPT_ESCAPE:
	case constant.DIALOG_TYPE_ACCEPT:
		if req.Next {
			args = append(args, lua.LTrue)
		} else {
			args = append(args, lua.LFalse)
		}
	}

	state, err := luax.Resume(co, args...)
	if err != nil {
		log.Fatal(err)
		return
	}

	if state == lua.ResumeYield {
		client.ch.Dialog = co
	}
}

func onGameClientNpcClick(ctx actor.Context, client *GameClientActor, req *req.NpcClick) {
	mapActor, ok := client.serverContext.MapActors[client.ch.Map]
	if !ok {
		return
	}

	ctx.Send(mapActor, &msg.SendMessage{
		OID: req.OID,
		Message: &msg.NpcClick{
			Sender: ctx.Self(),
		},
	})
}
