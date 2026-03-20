package wz

import "github.com/boyism80/fm/core/luax"

var (
	_ luax.Luable = (*ItemCore)(nil)
	_ luax.Luable = (*EquipmentCore)(nil)
	_ luax.Luable = (*Weapon)(nil)
	_ luax.Luable = (*Armor)(nil)
	_ luax.Luable = (*Consume)(nil)
	_ luax.Luable = (*Pet)(nil)
	_ luax.Luable = (*MiscItem)(nil)
	_ luax.Luable = (*CashItem)(nil)
	_ luax.Luable = (*Installation)(nil)
)
