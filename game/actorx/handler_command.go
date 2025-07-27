package actorx

import (
	"fmt"
	"log"
	"math"
	"path/filepath"
	"strconv"
	"strings"

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
	handler.RegisterCommandHandler("힌트", ctx, client, h, onCommandHint)
	handler.RegisterCommandHandler("공지", ctx, client, h, onCommandNotice)
	handler.RegisterCommandHandler("좌표", ctx, client, h, onCommandGetPosition)
	handler.RegisterCommandHandler("체력바꾸기", ctx, client, h, onCommandChangeHp)
	handler.RegisterCommandHandler("마력바꾸기", ctx, client, h, onCommandChangeMp)
	handler.RegisterCommandHandler("힘바꾸기", ctx, client, h, onCommandChangeStr)
	handler.RegisterCommandHandler("민첩바꾸기", ctx, client, h, onCommandChangeDex)
	handler.RegisterCommandHandler("지능바꾸기", ctx, client, h, onCommandChangeInt)
	handler.RegisterCommandHandler("행운바꾸기", ctx, client, h, onCommandChangeLuk)
	handler.RegisterCommandHandler("모든스탯바꾸기", ctx, client, h, onCommandChangeAllStats)
	handler.RegisterCommandHandler("레벨바꾸기", ctx, client, h, onCommandChangeLevel)
	handler.RegisterCommandHandler("무적", ctx, client, h, onCommandInvincible)
	handler.RegisterCommandHandler("직업바꾸기", ctx, client, h, onCommandChangeClass)
}

func onCommandCreateItem(ctx actor.Context, client *GameClientActor, params ...string) {
	if len(params) < 1 {
		log.Println("onCreateItem: missing itemId")
		return
	}

	itemId, err := strconv.Atoi(params[0])
	if err != nil {
		value, ok := client.serverContext.Resources.NameToItem(params[0])
		if !ok {
			log.Println("onCreateItem: invalid itemId:", params[0])
			return
		}
		itemId = int(value)
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
			constant.STAT_MESO: client.ch.Meso,
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
			constant.STAT_MESO: client.ch.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func onCommandFullMeso(ctx actor.Context, client *GameClientActor, params ...string) {
	client.ch.Meso = math.MaxInt32
	client.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: client.ch.Meso,
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
		mapId, ok := client.serverContext.Resources.NameToMap(params[0])
		if !ok {
			log.Println("onCommandChangeMap: invalid mapId:", params[0])
			return
		}
		value = int(mapId)
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
	animationType := constant.MOB_DIE_ANIMATION_TYPE_FADE_OUT
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

	mobId64, err := strconv.ParseUint(params[0], 10, 32)
	if err != nil {
		mobId, ok := client.serverContext.Resources.NameToMob(params[0])
		if !ok {
			log.Println("onCommandSpawnMob: invalid mobId:", params[0])
			return
		}
		mobId64 = uint64(mobId)
	}
	mobId := uint32(mobId64)

	mapActor, ok := client.serverContext.MapActors[client.ch.Map]
	if !ok {
		return
	}

	ctx.Send(mapActor, &msg.MapSpawningMob{
		MobId:    uint32(mobId),
		Position: client.ch.Position,
	})
}

func onCommandHint(ctx actor.Context, client *GameClientActor, params ...string) {
	if len(params) < 1 {
		log.Println("onCommandHint: missing text")
		return
	}

	text := params[0]
	width := 0
	if len(params) > 1 {
		if val, err := strconv.Atoi(params[1]); err == nil {
			width = val
		}
	}

	height := 0
	if len(params) > 2 {
		if val, err := strconv.Atoi(params[2]); err == nil {
			height = val
		}
	}

	client.Send(&resp.Hint{
		Text:   text,
		Width:  uint16(width),
		Height: uint16(height),
	}, types.SEND_POLICY_ENCRYPT)
}

func onCommandNotice(ctx actor.Context, client *GameClientActor, params ...string) {
	if len(params) < 1 {
		log.Println("onCommandNotice: missing text")
		return
	}
	text := params[0]

	var sb strings.Builder
	sb.WriteString(client.ch.Name)
	sb.WriteString(" : ")
	sb.WriteString(text)
	formattedText := sb.String()

	noticeType := resp.MSG_NOTICE
	if len(params) > 1 {
		if val, err := strconv.Atoi(params[1]); err == nil {
			noticeType = resp.ServerMessageType(val)
		}
	}

	megaEar := false
	if len(params) > 2 {
		if val, err := strconv.Atoi(params[2]); err == nil {
			megaEar = val != 0
		}
	}

	client.Send(&resp.Notice{
		Type:    noticeType,
		Channel: 0,
		Message: formattedText,
		MegaEar: megaEar,
	}, types.SEND_POLICY_ENCRYPT)
}

func changeStat(ctx actor.Context, client *GameClientActor, stats []constant.Stat, fn func(ch *entity.Character, value int) int, params ...string) {
	if len(params) < 1 {
		log.Println("onCommandChangeStr: missing str")
		return
	}
	value, err := strconv.Atoi(params[0])
	if err != nil {
		log.Println("onCommandChangeStr: invalid str:", params[0])
		return
	}
	if value > 32767 {
		value = 32767
	}
	value = fn(client.ch, value)

	updates := map[constant.Stat]int32{}
	for _, stat := range stats {
		updates[stat] = int32(value)
	}

	client.Send(&resp.UpdateStats{
		Stats: updates,
	}, types.SEND_POLICY_ENCRYPT)
}

func onCommandChangeHp(ctx actor.Context, client *GameClientActor, params ...string) {
	changeStat(ctx, client, []constant.Stat{constant.STAT_HP, constant.STAT_MAX_HP}, func(ch *entity.Character, value int) int {
		ch.Hp = uint16(value)
		ch.MaxHp = uint16(value)
		return value
	}, params...)
}

func onCommandChangeMp(ctx actor.Context, client *GameClientActor, params ...string) {
	changeStat(ctx, client, []constant.Stat{constant.STAT_MP, constant.STAT_MAX_MP}, func(ch *entity.Character, value int) int {
		ch.Mp = uint16(value)
		ch.MaxMp = uint16(value)
		return value
	}, params...)
}

func onCommandChangeStr(ctx actor.Context, client *GameClientActor, params ...string) {
	changeStat(ctx, client, []constant.Stat{constant.STAT_STR}, func(ch *entity.Character, value int) int {
		ch.Str = uint16(value)
		return value
	}, params...)
}

func onCommandChangeDex(ctx actor.Context, client *GameClientActor, params ...string) {
	changeStat(ctx, client, []constant.Stat{constant.STAT_DEX}, func(ch *entity.Character, value int) int {
		ch.Dex = uint16(value)
		return value
	}, params...)
}

func onCommandChangeInt(ctx actor.Context, client *GameClientActor, params ...string) {
	changeStat(ctx, client, []constant.Stat{constant.STAT_INT}, func(ch *entity.Character, value int) int {
		ch.Int = uint16(value)
		return value
	}, params...)
}

func onCommandChangeLuk(ctx actor.Context, client *GameClientActor, params ...string) {
	changeStat(ctx, client, []constant.Stat{constant.STAT_LUK}, func(ch *entity.Character, value int) int {
		ch.Luk = uint16(value)
		return value
	}, params...)
}

func onCommandChangeAllStats(ctx actor.Context, client *GameClientActor, params ...string) {
	changeStat(ctx, client, []constant.Stat{constant.STAT_STR, constant.STAT_DEX, constant.STAT_INT, constant.STAT_LUK}, func(ch *entity.Character, value int) int {
		ch.Str = uint16(value)
		ch.Dex = uint16(value)
		ch.Int = uint16(value)
		ch.Luk = uint16(value)
		return value
	}, params...)
}

func onCommandChangeLevel(ctx actor.Context, client *GameClientActor, params ...string) {
	changeStat(ctx, client, []constant.Stat{constant.STAT_LEVEL}, func(ch *entity.Character, value int) int {
		value = max(1, min(value, 255))
		ch.Level = uint8(value)
		return value
	}, params...)
}

func onCommandInvincible(ctx actor.Context, client *GameClientActor, params ...string) {
	if len(params) < 1 {
		log.Println("onCommandInvincible: missing invincible")
		return
	}
	invincible, err := strconv.Atoi(params[0])
	if err != nil {
		log.Println("onCommandInvincible: invalid invincible:", params[0])
		return
	}
	client.ch.Invincible = invincible != 0
	client.ch.Listener.OnMessage(fmt.Sprintf("무적 상태 변경: %v", client.ch.Invincible))
}

func onCommandChangeClass(ctx actor.Context, client *GameClientActor, params ...string) {
	nameToClassId := map[string]int{
		"초보자":      0,
		"전사":       100,
		"파이터":      110,
		"크루세이더":    111,
		"히어로":      112,
		"페이지":      120,
		"나이트":      121,
		"스피어맨":     130,
		"용기사":      131,
		"다크나이트":    132,
		"마법사":      200,
		"썬콜워저드":    210,
		"썬콜매이지":    211,
		"썬콜아크메이지":  212,
		"불독위저드":    220,
		"불독메이지":    221,
		"불독아크메이지":  222,
		"클레릭":      230,
		"비숍":       232,
		"궁수":       300,
		"헌터":       310,
		"레인져":      311,
		"보우마스터":    312,
		"사수":       320,
		"저격수":      321,
		"신궁":       322,
		"도적":       400,
		"어쌔신":      410,
		"허밋":       411,
		"나이트로드":    412,
		"시프":       420,
		"시프마스터":    421,
		"섀도어":      422,
		"해적":       500,
		"인파이터":     510,
		"버커니어":     511,
		"바이퍼":      512,
		"건슬링거":     520,
		"발키리":      521,
		"캡틴":       522,
		"노블레스":     1000,
		"소울마스터":    1100,
		"소울마스터2차":  1110,
		"소울마스터3차":  1111,
		"플레임위드":    1200,
		"플레임위드2차":  1210,
		"플레임위드3차":  1211,
		"윈드브래이커":   1300,
		"윈드브레이커2차": 1310,
		"윈드브레이커3차": 1311,
		"나이트워커":    1400,
		"나이트워커2차":  1410,
		"나이트워커3차":  1411,
		"스트라이커":    1500,
		"스트라이커2차":  1510,
		"스트라이커3차":  1511,
	}

	if len(params) < 1 {
		log.Println("onCommandChangeClass: missing class")
		return
	}

	classId, ok := nameToClassId[params[0]]
	if !ok {
		log.Println("onCommandChangeClass: invalid class:", params[0])
		return
	}

	client.ch.Class = uint16(classId)
	client.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_JOB: int32(classId),
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func onCommandGetPosition(ctx actor.Context, client *GameClientActor, params ...string) {
	client.listener.OnChat(fmt.Sprintf("%d, %d, %d", client.ch.Map, client.ch.Position.X, client.ch.Position.Y), false, false)
}
