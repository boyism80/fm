package constant

// ServerMessageType defines the different SERVERMESSAGE types.
type ServerMessageType uint8

const (
	MSG_NOTICE                ServerMessageType = 0  // [Notice]
	MSG_POPUP                 ServerMessageType = 1  // Popup
	MSG_MEGAPHONE             ServerMessageType = 2  // Megaphone
	MSG_SUPER_MEGAPHONE       ServerMessageType = 3  // Super Megaphone
	MSG_SCROLLING_TOP         ServerMessageType = 4  // Scrolling message at top
	MSG_PINK_TEXT             ServerMessageType = 5  // Pink Text
	MSG_LIGHT_BLUE_TEXT       ServerMessageType = 6  // Lightblue Text
	MSG_ITEM_MEGAPHONE        ServerMessageType = 8  // Item megaphone
	MSG_HEART_MEGAPHONE       ServerMessageType = 9  // Heart megaphone
	MSG_SKULL_SUPER_MEGAPHONE ServerMessageType = 10 // Skull Super megaphone
	MSG_GREEN_MEGAPHONE       ServerMessageType = 11 // Green megaphone message
	MSG_THREE_MEGAPHONE_LINES ServerMessageType = 12 // Three-line megaphone text
	MSG_EOF                   ServerMessageType = 13 // End of file
	MSG_ANI                   ServerMessageType = 14 // Ani msg
	MSG_RED_GACHAPON_BOX      ServerMessageType = 15 // Red Gachapon box
	MSG_BLUE_NOTICE           ServerMessageType = 18 // Blue Notice
)
