package constant

// DialogType represents the type of NPC dialog interaction.
type DialogType uint8

const (
	DIALOG_TYPE_DEFAULT       DialogType = 0  // Simple message dialog
	DIALOG_TYPE_YES_NO        DialogType = 1  // Yes/No confirmation
	DIALOG_TYPE_INPUT         DialogType = 2  // Text input dialog
	DIALOG_TYPE_LIST          DialogType = 4  // Selection list
	DIALOG_TYPE_ACCEPT_ESCAPE DialogType = 11 // Accept or escape
	DIALOG_TYPE_ACCEPT        DialogType = 12 // Accept only
)
