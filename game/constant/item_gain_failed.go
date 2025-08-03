package constant

// ItemGainFailedType defines the different item gain failure types.
type ItemGainFailedType uint8

const (
	ITEM_GAIN_FAILED_TYPE_FULL  ItemGainFailedType = 0xFF
	ITEM_GAIN_FAILED_TYPE_ERROR ItemGainFailedType = 0xFE
)
