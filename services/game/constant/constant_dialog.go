package constant

type DialogType uint8

const (
	DIALOG_TYPE_DEFAULT       DialogType = 0
	DIALOG_TYPE_YES_NO        DialogType = 1
	DIALOG_TYPE_INPUT         DialogType = 2
	DIALOG_TYPE_LIST          DialogType = 4
	DIALOG_TYPE_ACCEPT_ESCAPE DialogType = 11
	DIALOG_TYPE_ACCEPT        DialogType = 12
)
