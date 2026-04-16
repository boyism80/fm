package entity

import "github.com/boyism80/fm/protocol/dto"

type KeyLayout struct {
	slots map[int]keyBinding
}

type keyBinding struct {
	typ    byte
	action int32
}

func NewKeyLayout() *KeyLayout {
	return &KeyLayout{slots: make(map[int]keyBinding)}
}

func (k *KeyLayout) IsEmpty() bool {
	return k == nil || len(k.slots) == 0
}

func (k *KeyLayout) ApplyChange(slot int, typ byte, action int32) {
	if k == nil {
		return
	}
	if typ != 0 {
		k.slots[slot] = keyBinding{typ: typ, action: action}
		return
	}
	delete(k.slots, slot)
}

func (k *KeyLayout) Bindings() map[int]dto.KeyBinding {
	if k == nil || len(k.slots) == 0 {
		return map[int]dto.KeyBinding{}
	}
	out := make(map[int]dto.KeyBinding, len(k.slots))
	for slot, b := range k.slots {
		out[slot] = dto.KeyBinding{
			Type:   b.typ,
			Action: b.action,
		}
	}
	return out
}
