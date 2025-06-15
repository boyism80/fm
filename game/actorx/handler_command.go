package actorx

import (
	"log"
	"math"
	"path/filepath"
	"strconv"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/common/luax"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol/resp"
)

func RegisterGameClientCommandHandler(ctx actor.Context, client *GameClientActor, h *handler.CommandHandler) {
	handler.RegisterCommandHandler("아이템", ctx, client, h, onCommandCreateItem)
	handler.RegisterCommandHandler("메소초기화", ctx, client, h, onCommandClearMeso)
	handler.RegisterCommandHandler("메소얻기", ctx, client, h, onCommandGainMeso)
	handler.RegisterCommandHandler("풀메소", ctx, client, h, onCommandFullMeso)
	handler.RegisterCommandHandler("맵이동", ctx, client, h, onCommandChangeMap)
	handler.RegisterCommandHandler("다이얼로그", ctx, client, h, onCommandDialog)
	handler.RegisterCommandHandler("스크립트", ctx, client, h, onCommandScript)
	handler.RegisterCommandHandler("몬스터죽이기", ctx, client, h, onCommandMobKill)
	handler.RegisterCommandHandler("몬스터생성", ctx, client, h, onCommandSpawnMob)
}

func onCommandCreateItem(ctx actor.Context, client *GameClientActor, params ...string) {
	if len(params) < 1 {
		log.Println("onCreateItem: missing itemId")
		return
	}

	itemId, err := strconv.Atoi(params[0])
	if err != nil {
		log.Println("onCreateItem: invalid itemId:", params[0])
		return
	}

	count := 1
	if len(params) >= 2 {
		if parsedCount, err := strconv.Atoi(params[1]); err == nil {
			count = parsedCount
		} else {
			log.Println("onCreateItem: invalid count, using default 1:", params[1])
		}
	}

	item, err := entity.NewItem(client.serverContext, uint32(itemId), uint16(count))
	if err != nil {
		log.Println(err)
		return
	}

	inventoryType := item.GetInventoryType()
	inventory := client.ch.Inventory[inventoryType]
	nextSlot, ok := inventory.NextSlot()
	if !ok {
		return
	}

	inventory.Items[int16(nextSlot)] = item
	client.Send(&resp.AddItem{
		IsDrop:        false,
		Slot:          nextSlot,
		InventoryType: inventoryType,
		Item:          item,
	}, types.SEND_POLICY_ENCRYPT)
}

func onCommandClearMeso(ctx actor.Context, client *GameClientActor, params ...string) {
	client.ch.Meso = 0
	client.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.StatMeso: client.ch.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func onCommandGainMeso(ctx actor.Context, client *GameClientActor, params ...string) {
	if len(params) < 1 {
		log.Println("onGainMeso: missing itemId")
		return
	}

	value, err := strconv.Atoi(params[0])
	if err != nil {
		log.Println("onGainMeso: invalid itemId:", params[0])
		return
	}
	cap := math.MaxInt32 - client.ch.Meso
	meso := min(cap, int32(value))
	client.ch.Meso += meso
	client.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.StatMeso: client.ch.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func onCommandFullMeso(ctx actor.Context, client *GameClientActor, params ...string) {
	client.ch.Meso = math.MaxInt32
	client.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.StatMeso: client.ch.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func onCommandChangeMap(ctx actor.Context, client *GameClientActor, params ...string) {
	if len(params) < 1 {
		log.Println("onGainMeso: missing itemId")
		return
	}

	value, err := strconv.Atoi(params[0])
	if err != nil {
		log.Println("onGainMeso: invalid itemId:", params[0])
		return
	}

	client.Warp(ctx, uint32(value), 0)
}

func onCommandDialog(ctx actor.Context, client *GameClientActor, params ...string) {
	prev := 0
	if len(params) > 1 {
		if val, err := strconv.Atoi(params[1]); err == nil {
			prev = val
		}
	}

	next := 0
	if len(params) > 2 {
		if val, err := strconv.Atoi(params[2]); err == nil {
			next = val
		}
	}

	if client.ch.Listener != nil {
		client.ch.Listener.OnDialog(9001000, "안녕하세요", prev != 0, next != 0)
	}
}

func onCommandScript(ctx actor.Context, client *GameClientActor, params ...string) {
	fileName := "script.lua"
	if len(params) >= 1 {
		fileName = params[0]
	}
	path := filepath.Join("script", fileName)

	luax.Call(ctx, ctx.Self(), path, "on_start", client.builtin)
}

func onCommandMobKill(ctx actor.Context, client *GameClientActor, params ...string) {
	animationType := constant.MobDieAnimationTypeFadeOut
	if len(params) > 0 {
		v, err := strconv.Atoi(params[0])
		if err != nil {
			log.Println("onMobKill: invalid animationType:", params[0])
			return
		}
		animationType = constant.MobDieAnimationType(v)
	}

	mapActor, ok := client.serverContext.MapActors[client.ch.Map]
	if !ok {
		return
	}

	ctx.Send(mapActor, &msg.MapClearMobs{
		AnimationType: animationType,
	})
}

func onCommandSpawnMob(ctx actor.Context, client *GameClientActor, params ...string) {

	if len(params) < 1 {
		log.Println("onCommandSpawnMob: missing mobId")
		return
	}

	mobId, err := strconv.Atoi(params[0])
	if err != nil {
		log.Println("onCommandSpawnMob: invalid mobId:", params[0])
		return
	}

	mapActor, ok := client.serverContext.MapActors[client.ch.Map]
	if !ok {
		return
	}

	ctx.Send(mapActor, &msg.MapSpawningMob{
		MobId:    uint32(mobId),
		Position: client.ch.Position,
	})
}
