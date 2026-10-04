package bot

import (
	"errors"
	"fmt"
	"time"

	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/bot/conn"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

var ErrConnect = errors.New("connect")

type outbound interface {
	Opcode() byte
	Serialize(writer *stream.StreamWriter) error
}

type Bot struct {
	Index  int
	ID     string
	Name   string
	Gender uint8
	CharID uint32
	Gen    int
	Map    uint32
	Spawn  uint8
	Moved  *types.Point[int16]
	HP     uint16
	Level  uint8
	Job    uint16
	EXP    int32
	Meso   int32
	Fame   int32
	NPCs   map[uint32]uint32
	Quests map[uint16]uint8
	Items  map[constant.InventoryType]map[int16]ItemSlot
	Dialog constant.DialogType
	cfg    *config.Bot
	conn   *conn.Conn
}

type ItemSlot struct {
	ItemID uint32
	Count  uint16
}

func New(cfg *config.Bot, runID string, n int) (*Bot, error) {
	name := fmt.Sprintf("b%s%d", runID, n)
	if len(name) > 12 {
		return nil, fmt.Errorf("character name %q is longer than 12 bytes", name)
	}
	return &Bot{
		Index:  n,
		ID:     fmt.Sprintf("bot_%s_%d", runID, n),
		Name:   name,
		NPCs:   make(map[uint32]uint32),
		Quests: make(map[uint16]uint8),
		Items:  make(map[constant.InventoryType]map[int16]ItemSlot),
		cfg:    cfg,
	}, nil
}

func (b *Bot) Send(pkt outbound) error {
	if b.conn == nil {
		return fmt.Errorf("bot %s is not connected", b.Name)
	}
	return b.conn.Send(pkt)
}

func (b *Bot) Close() {
	if b.conn == nil {
		return
	}
	_ = b.conn.Close()
	b.conn = nil
}

func (b *Bot) Listen(onPacket func(gen int, opcode uint16, body []byte), onClose func(gen int, err error)) {
	c, gen := b.conn, b.Gen
	go func() {
		for {
			opcode, body, err := c.Read()
			if err != nil {
				onClose(gen, err)
				return
			}
			onPacket(gen, opcode, body)
		}
	}()
}

func (b *Bot) Update(pkt any) {
	switch p := pkt.(type) {
	case *response.Login:
		b.Moved = nil
		clear(b.NPCs)
		clear(b.Quests)
		clear(b.Items)
		if p.Character == nil {
			return
		}
		b.Map, b.Spawn, b.HP = p.Character.Map, p.Character.SpawnPoint, p.Character.Hp
		b.Level, b.Job, b.EXP, b.Fame = p.Character.Level, p.Character.Class, int32(p.Character.Exp), int32(p.Character.Population)
		b.Meso = p.Character.Inventory.Meso
		for typ, tab := range p.Character.Inventory.Tabs {
			for slot, item := range tab.Items {
				b.setItem(typ, slot, item)
			}
		}
		for parts, item := range p.Character.Inventory.Equipped {
			b.setItem(constant.InventoryTypeEquipment, int16(parts), item)
		}
		for _, q := range p.Character.QuestsStarted {
			b.Quests[q.QuestID] = q.Status
		}
		for _, q := range p.Character.QuestsCompleted {
			b.Quests[q.QuestID] = q.Status
		}
	case *response.Warp:
		if p.Character != nil {
			b.Map, b.Spawn, b.HP = p.Character.Map, p.Character.SpawnPoint, p.Character.Hp
		}
		b.Moved = nil
		clear(b.NPCs)
	case *response.UpdateStats:
		for stat, value := range p.Stats {
			switch stat {
			case constant.StatHP:
				b.HP = uint16(value)
			case constant.StatLevel:
				b.Level = uint8(value)
			case constant.StatClass:
				b.Job = uint16(value)
			case constant.StatEXP:
				b.EXP = value
			case constant.StatMeso:
				b.Meso = value
			case constant.StatPopulation:
				b.Fame = value
			}
		}
	case *response.UpdateQuest:
		if p.Status == response.QuestWireStatusNotStarted {
			delete(b.Quests, p.QuestID)
		} else {
			b.Quests[p.QuestID] = p.Status
		}
	case *response.InventoryOperation:
		for _, change := range p.Changes {
			tab := b.Items[change.InventoryType]
			switch change.Mode {
			case response.INVENTORY_MODE_ADD:
				b.setItem(change.InventoryType, change.Slot, change.Item)
			case response.INVENTORY_MODE_UPDATE:
				if slot, ok := tab[change.Slot]; ok {
					slot.Count = change.Count
					tab[change.Slot] = slot
				}
			case response.INVENTORY_MODE_MOVE:
				src, ok := tab[change.Slot]
				if ok == false {
					continue
				}
				if dest, ok := tab[change.Dest]; ok {
					tab[change.Slot] = dest
				} else {
					delete(tab, change.Slot)
				}
				tab[change.Dest] = src
			case response.INVENTORY_MODE_REMOVE:
				delete(tab, change.Slot)
			}
		}
	case *response.SpawnNpc:
		b.NPCs[p.NPC.OID] = p.NPC.NpcId
	case *response.RemoveNpc:
		delete(b.NPCs, p.OID)
	case *response.Dialog:
		b.Dialog = p.Type
	case *response.DialogYesNo:
		b.Dialog = constant.DialogTypeYesNo
	case *response.DialogInput:
		b.Dialog = constant.DialogTypeInput
	case *response.DialogList:
		b.Dialog = constant.DialogTypeList
	case *response.DialogStyle:
		b.Dialog = constant.DialogTypeStyle
	case *response.DialogAccept:
		b.Dialog = constant.DialogTypeAccept
		if p.EnableEscape {
			b.Dialog = constant.DialogTypeAcceptEscape
		}
	}
}

func (b *Bot) setItem(typ constant.InventoryType, slot int16, item dto.Item) {
	if item == nil {
		return
	}
	tab, ok := b.Items[typ]
	if ok == false {
		tab = make(map[int16]ItemSlot)
		b.Items[typ] = tab
	}
	tab[slot] = ItemSlot{ItemID: item.GetItemID(), Count: item.GetCount()}
}

func (b *Bot) FindNPC(templateID uint32) (uint32, bool) {
	for oid, id := range b.NPCs {
		if id == templateID {
			return oid, true
		}
	}
	return 0, false
}

func (b *Bot) Move(pos types.Point[int16], foothold int16) error {
	err := b.Send(&request.MovePlayer{Fragments: []dto.MoveFragment{
		&dto.AbsoluteLifeMovement{BasicMovement: &dto.BasicMovement{}, Position: pos, Foothold: foothold},
	}})
	if err != nil {
		return err
	}
	b.Moved = &pos
	return nil
}

func (b *Bot) Position(resources *wz.Resources) (types.Point[int16], bool) {
	if b.Moved != nil {
		return *b.Moved, true
	}
	m, ok := resources.Maps[b.Map]
	if ok == false {
		return types.Point[int16]{}, false
	}
	portal, ok := m.Portals[b.Spawn]
	if ok == false {
		return types.Point[int16]{}, false
	}
	return portal.Position, true
}

func (b *Bot) wait() time.Duration {
	return time.Duration(b.cfg.TimeoutMs) * time.Millisecond
}
