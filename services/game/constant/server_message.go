package constant

type ServerMessageType uint8

const (
	MSG_NOTICE                ServerMessageType = 0
	MSG_POPUP                 ServerMessageType = 1
	MSG_MEGAPHONE             ServerMessageType = 2
	MSG_SUPER_MEGAPHONE       ServerMessageType = 3
	MSG_SCROLLING_TOP         ServerMessageType = 4
	MSG_PINK_TEXT             ServerMessageType = 5
	MSG_LIGHT_BLUE_TEXT       ServerMessageType = 6
	MSG_ITEM_MEGAPHONE        ServerMessageType = 8
	MSG_HEART_MEGAPHONE       ServerMessageType = 9
	MSG_SKULL_SUPER_MEGAPHONE ServerMessageType = 10
	MSG_GREEN_MEGAPHONE       ServerMessageType = 11
	MSG_THREE_MEGAPHONE_LINES ServerMessageType = 12
	MSG_EOF                   ServerMessageType = 13
	MSG_ANI                   ServerMessageType = 14
	MSG_RED_GACHAPON_BOX      ServerMessageType = 15
	MSG_BLUE_NOTICE           ServerMessageType = 18
)
