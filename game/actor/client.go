package actor

import (
	"log"
	"net"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/crypt"
	"github.com/boyism80/fm/common/handler"
	common_msg "github.com/boyism80/fm/common/msg"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/common/util"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/data"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol/resp"
)

type GameClientActor struct {
	ctx            *context.ServerContext
	ch             *entity.Character
	conn           net.Conn
	buffer         []byte
	messageHandler *handler.MessageHandler
	packetHandler  *handler.PacketHandler
	commandHandler *handler.CommandHandler
	sendCrypt      crypt.Encryption
	receiveCrypt   crypt.Encryption
	stopTimer      scheduler.CancelFunc
}

func NewGameClientActor(ctx actor.Context, serverCtx *context.ServerContext, conn net.Conn) actor.Actor {
	ivSend := []byte{0x2F, 0xA3, 0x65, 0x43}
	ivRecv := []byte{0x65, 0x56, 0x12, 0xFD}

	actor := &GameClientActor{
		ctx:            serverCtx,
		conn:           conn,
		buffer:         []byte{},
		messageHandler: handler.NewMessageHandler(),
		packetHandler:  handler.NewPacketHandler(),
		commandHandler: handler.NewCommandHandler(),
		sendCrypt:      crypt.NewEncryption(ivSend, -5),
		receiveCrypt:   crypt.NewEncryption(ivRecv, 5),
	}

	RegisterGameClientMessageHandlers(ctx, actor, actor.messageHandler)
	RegisterGameClientPacketHandler(ctx, actor, actor.packetHandler)
	RegisterGameClientCommandHandler(ctx, actor, actor.commandHandler)
	return actor
}

func (state *GameClientActor) Receive(context actor.Context) {
	state.messageHandler.Handle(context)
}

func (state *GameClientActor) Invoke(ctx actor.Context, header int, data []byte) error {
	return state.packetHandler.Handle(ctx, header, data)
}

func (state *GameClientActor) Send(p types.Packet, policy types.SendPolicy) {
	writer := stream.NewStreamWriter(stream.LittleEndian)
	p.Serialize(writer)
	bytes := writer.Bytes()

	log.Println("[S] " + util.ToHexString(bytes))

	if policy == types.SEND_POLICY_RAW {
		state.conn.Write(bytes)
		return
	}

	writer = stream.NewStreamWriter(stream.LittleEndian)
	if policy&types.SEND_POLICY_ENCRYPT != 0 {
		header := state.sendCrypt.GetPacketHeader(len(bytes))
		writer.Write(header)
		bytes = state.sendCrypt.Encrypt(bytes)
	}

	writer.Write(bytes)
	bytes = writer.Bytes()
	state.conn.Write(bytes)
}

func (client *GameClientActor) Drop(ctx actor.Context, invenType constant.InventoryType, slot int16, count uint16) {
	ch := client.ch
	item, ok := ch.Inventory[invenType].Items[slot]
	if !ok {
		return
	}

	mapActor := client.ctx.MapActors[ch.Map]
	if mapActor == nil {
		return
	}

	removed := (item.Reduce(count) == 0)
	if removed {
		client.Send(&resp.RemoveInventorySlot{
			InventoryType: invenType,
			Slot:          slot,
		}, types.SEND_POLICY_ENCRYPT)
	} else {
		client.Send(&resp.UpdateInventorySlot{
			InventoryType: invenType,
			Slot:          slot,
			Item:          item,
		}, types.SEND_POLICY_ENCRYPT)
	}

	spawned := item.Clone(count)
	spawned.BindDrop(&entity.Drop{
		Object:       &entity.Object{},
		Owner:        ch.ID,
		SpawnedPoint: ch.Position,
		DropType:     constant.DropTypeFFA,
		Looting:      false,
	})

	ctx.Send(mapActor, &msg.MapSpawnItem{
		Item:    spawned,
		Owner:   ctx.Self(),
		OwnerID: ch.ID,
	})

	if removed {
		delete(ch.Inventory[invenType].Items, slot)
	}
}

func (client *GameClientActor) MoveItem(invenType constant.InventoryType, sourceSlot int16, destSlot int16) {
	ch := client.ch
	inven := ch.Inventory[invenType]
	src, ok := inven.Items[sourceSlot]
	if !ok {
		return
	}

	dst, ok := inven.Items[destSlot]

	if !ok {
		inven.Items[destSlot] = inven.Items[sourceSlot]
		delete(inven.Items, sourceSlot)
		client.Send(&resp.SwapInventorySlot{
			InventoryType: invenType,
			Source:        sourceSlot,
			Dest:          destSlot,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}

	specSrc := src.GetSpec()
	specDst := dst.GetSpec()
	if specSrc != specDst {
		inven.Items[sourceSlot], inven.Items[destSlot] = inven.Items[destSlot], inven.Items[sourceSlot]
		client.Send(&resp.SwapInventorySlot{
			InventoryType: invenType,
			Source:        sourceSlot,
			Dest:          destSlot,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}

	limit := min(src.GetCount(), specSrc.GetCapacity()-dst.GetCount())
	dst.Increase(limit)
	if src.Reduce(limit) == 0 {
		client.Send(&resp.FullMergeInventorySlot{
			InventoryType: invenType,
			Source:        sourceSlot,
			Dest:          destSlot,
			Count:         dst.GetCount(),
		}, types.SEND_POLICY_ENCRYPT)
		delete(inven.Items, sourceSlot)
	} else {
		client.Send(&resp.PartialMergeInventorySlot{
			InventoryType: invenType,
			Source:        sourceSlot,
			Dest:          destSlot,
			SourceCount:   src.GetCount(),
			DestCount:     dst.GetCount(),
		}, types.SEND_POLICY_ENCRYPT)
	}
}

func (client *GameClientActor) Unequip(ctx actor.Context, parts constant.EquipmentPartsType, slot int16) {
	ch := client.ch
	if ch.Equipments[parts] == nil {
		return
	}

	inven := ch.Inventory[constant.InventoryTypeEquipment]
	if inven.Items[slot] != nil {
		return
	}

	mapActor := client.ctx.MapActors[client.ch.Map]
	if mapActor == nil {
		return
	}

	inven.Items[slot] = ch.Equipments[parts]
	delete(ch.Equipments, parts)
	client.Send(&resp.SwapInventorySlot{
		InventoryType:   constant.InventoryTypeEquipment,
		Source:          int16(parts),
		Dest:            slot,
		EquipmentAction: resp.EquipmentActionTypeOff,
	}, types.SEND_POLICY_ENCRYPT)

	ctx.Send(mapActor, &msg.MapBroadcastRange{
		Sender: ctx.Self(),
		Pivot:  ch.Position,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.UpdateCharacterLook{
				Character: ch,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
		ExceptSelf: true,
	})
}

func (client *GameClientActor) Equip(ctx actor.Context, parts constant.EquipmentPartsType, slot int16) {
	ch := client.ch
	inven := ch.Inventory[constant.InventoryTypeEquipment]
	if inven.Items[slot] == nil {
		return
	}

	mapActor := client.ctx.MapActors[client.ch.Map]
	if mapActor == nil {
		return
	}

	new, ok := inven.Items[slot].(*entity.Equipment)
	if !ok {
		return
	}
	old, swap := ch.Equipments[parts]
	switch parts {
	case constant.EquipmentPartsTop:
		if new.IsOverall() {
			_, isWearPants := ch.Equipments[constant.EquipmentPartsPants]
			storageSlot, isFree := inven.NextSlot()
			if isWearPants {
				if !isFree {
					client.Send(&resp.ItemGainFailed{
						Mode: resp.ItemGainFailedTypeFull,
					}, types.SEND_POLICY_ENCRYPT)
					return
				}
				client.Unequip(ctx, constant.EquipmentPartsPants, int16(storageSlot))
			}
		}

	case constant.EquipmentPartsPants:
		top, isWearTop := ch.Equipments[constant.EquipmentPartsTop]
		if isWearTop && top.IsOverall() {
			storageSlot, isFree := inven.NextSlot()
			if swap && !isFree {
				client.Send(&resp.ItemGainFailed{
					Mode: resp.ItemGainFailedTypeFull,
				}, types.SEND_POLICY_ENCRYPT)
				return
			}

			client.Unequip(ctx, constant.EquipmentPartsTop, int16(storageSlot))
		}
	}
	ch.Equipments[parts], inven.Items[slot] = new, old
	if !swap {
		delete(inven.Items, slot)
	}

	client.Send(&resp.SwapInventorySlot{
		InventoryType:   constant.InventoryTypeEquipment,
		Source:          slot,
		Dest:            int16(parts),
		EquipmentAction: resp.EquipmentActionTypeOn,
	}, types.SEND_POLICY_ENCRYPT)

	ctx.Send(mapActor, &msg.MapBroadcastRange{
		Sender: ctx.Self(),
		Pivot:  ch.Position,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.UpdateCharacterLook{
				Character: ch,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
		ExceptSelf: true,
	})
}

func (client *GameClientActor) MergeItems(inventoryType constant.InventoryType) {
	inven := client.ch.Inventory[inventoryType]
	buckets := map[data.ItemSpec]map[int16]entity.Item{}

	for i := range inven.SlotLimit {
		item := inven.Items[int16(i+1)]
		if item == nil {
			continue
		}

		spec := item.GetSpec()
		if buckets[spec] == nil {
			buckets[spec] = map[int16]entity.Item{}
		}

		buckets[spec][int16(i+1)] = item
	}

	for spec, bucket := range buckets {
		count := uint16(0)
		for _, v := range bucket {
			count += v.GetCount()
		}

		capacity := spec.GetCapacity()
		for slot, item := range bucket {
			value := min(capacity, count)
			if item.GetCount() != value {
				item.SetCount(value)
				if value == 0 {
					client.Send(&resp.RemoveInventorySlot{
						InventoryType: inventoryType,
						Slot:          slot,
					}, types.SEND_POLICY_ENCRYPT)
					delete(inven.Items, slot)
				} else {
					client.Send(&resp.UpdateInventorySlot{
						InventoryType: inventoryType,
						Slot:          slot,
						Item:          item,
					}, types.SEND_POLICY_ENCRYPT)
				}
			}
			count -= value
		}
	}
}

func (client *GameClientActor) SortInventory(inventoryType constant.InventoryType) {
	inven := client.ch.Inventory[inventoryType]
	n := inven.SlotLimit
	buffer := make([]entity.Item, n)
	for i := range n {
		buffer[i] = inven.Items[int16(i+1)]
	}

	less := func(item1, item2 entity.Item) bool {
		if item1 == nil && item2 == nil {
			return false
		}
		if item1 == nil {
			return false
		}
		if item2 == nil {
			return true
		}
		id1, id2 := item1.GetSpec().GetID(), item2.GetSpec().GetID()
		if id1 != id2 {
			return id1 < id2
		}
		return item1.GetCount() > item2.GetCount()
	}

	partition := func(low, high int) int {
		pivot := buffer[(low+high)/2]
		i1, i2 := low, high
		for i1 <= i2 {
			for less(buffer[i1], pivot) {
				i1++
			}
			for less(pivot, buffer[i2]) {
				i2--
			}
			if i1 <= i2 {

				buffer[i1], buffer[i2] = buffer[i2], buffer[i1]

				client.Send(&resp.SwapInventorySlot{
					InventoryType: inventoryType,
					Source:        int16(i1 + 1),
					Dest:          int16(i2 + 1),
				}, types.SEND_POLICY_ENCRYPT)
				i1++
				i2--
			}
		}
		return i1
	}

	var qsort func(low, high int)
	qsort = func(low, high int) {
		if low < high {
			p := partition(low, high)
			qsort(low, p-1)
			qsort(p, high)
		}
	}

	qsort(0, int(n-1))

	inven.Items = map[int16]entity.Item{}
	for i := range buffer {
		if buffer[i] != nil {
			inven.Items[int16(i+1)] = buffer[i]
		}
	}
}
