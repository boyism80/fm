package constant

type DialogType uint8

const (
	DialogTypeDefault      DialogType = 0
	DialogTypeYesNo        DialogType = 1
	DialogTypeInput        DialogType = 2
	DialogTypeList         DialogType = 4
	DialogTypeAcceptEscape DialogType = 11
	DialogTypeAccept       DialogType = 12
)
