package constant

type ItemGainFailedType uint8

const (
	ItemGainFailedTypeFull  ItemGainFailedType = 0xFF
	ItemGainFailedTypeError ItemGainFailedType = 0xFE
)

type ShowItemGainType uint8

const (
	ShowItemGainTypeStatus ShowItemGainType = iota
	ShowItemGainTypeChat
)

type ShowMesoGainType uint8

const (
	ShowMesoGainTypeStatus ShowMesoGainType = iota
	ShowMesoGainTypeChat
)
