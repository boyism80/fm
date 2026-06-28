package constant

type ServerMessageType uint8

const (
	MsgNotice              ServerMessageType = 0
	MsgPopup               ServerMessageType = 1
	MsgMegaphone           ServerMessageType = 2
	MsgSuperMegaphone      ServerMessageType = 3
	MsgScrollingTop        ServerMessageType = 4
	MsgPinkText            ServerMessageType = 5
	MsgLightBlueText       ServerMessageType = 6
	MsgItemMegaphone       ServerMessageType = 8
	MsgHeartMegaphone      ServerMessageType = 9
	MsgSkullSuperMegaphone ServerMessageType = 10
	MsgGreenMegaphone      ServerMessageType = 11
	MsgThreeMegaphoneLines ServerMessageType = 12
	MsgEOF                 ServerMessageType = 13
	MsgANI                 ServerMessageType = 14
	MsgRedGachaponBox      ServerMessageType = 15
	MsgBlueNotice          ServerMessageType = 18
)

func AllServerMessageTypes() map[string]ServerMessageType {
	return map[string]ServerMessageType{
		"Notice":              MsgNotice,
		"Popup":               MsgPopup,
		"Megaphone":           MsgMegaphone,
		"SuperMegaphone":      MsgSuperMegaphone,
		"ScrollingTop":        MsgScrollingTop,
		"PinkText":            MsgPinkText,
		"LightBlueText":       MsgLightBlueText,
		"ItemMegaphone":       MsgItemMegaphone,
		"HeartMegaphone":      MsgHeartMegaphone,
		"SkullSuperMegaphone": MsgSkullSuperMegaphone,
		"GreenMegaphone":      MsgGreenMegaphone,
		"ThreeMegaphoneLines": MsgThreeMegaphoneLines,
		"EOF":                 MsgEOF,
		"ANI":                 MsgANI,
		"RedGachaponBox":      MsgRedGachaponBox,
		"BlueNotice":          MsgBlueNotice,
	}
}

const DoorNoTownPortalMessage = "마을의 미스틱 도어 지점이 꽉 차서 지금은 사용할 수 없습니다."
