package entity

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

func (k *KeyLayout) ToGrpcDTO() []*internal.KeyLayoutBinding {
	if k == nil || len(k.slots) == 0 {
		return nil
	}
	out := make([]*internal.KeyLayoutBinding, 0, len(k.slots))
	for slot, b := range k.slots {
		out = append(out, &internal.KeyLayoutBinding{
			Slot:   int32(slot),
			Type:   uint32(b.typ),
			Action: b.action,
		})
	}
	return out
}

func (k *KeyLayout) LoadKeyLayoutProto(list []*internal.KeyLayoutBinding) {
	if k == nil {
		return
	}
	k.slots = make(map[int]keyBinding)
	for _, b := range list {
		if b == nil {
			continue
		}
		k.slots[int(b.GetSlot())] = keyBinding{
			typ:    byte(b.GetType()),
			action: b.GetAction(),
		}
	}
}
