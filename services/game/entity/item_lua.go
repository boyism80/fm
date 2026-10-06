package entity

import "github.com/boyism80/fm/core/luax"

var (
	_ luax.Luable = (*ItemCore)(nil)
	_ luax.Luable = (*EquipmentCore)(nil)
	_ luax.Luable = (*FieldPlacement)(nil)
	_ luax.Luable = (*Meso)(nil)
	_ luax.Luable = (*Consume)(nil)
	_ luax.Luable = (*CashItem)(nil)
	_ luax.Luable = (*MiscItem)(nil)
	_ luax.Luable = (*Installation)(nil)
	_ luax.Luable = (*Pet)(nil)
)
