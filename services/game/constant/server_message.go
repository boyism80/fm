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

type MessageScope int

const (
	MessageScopeSelf MessageScope = iota
	MessageScopeMap
	MessageScopeChannel
	MessageScopeWorld
)

func AllMessageScopes() map[string]MessageScope {
	return map[string]MessageScope{
		"Self":    MessageScopeSelf,
		"Map":     MessageScopeMap,
		"Channel": MessageScopeChannel,
		"World":   MessageScopeWorld,
	}
}

const DoorNoTownPortalMessage = "마을의 미스틱 도어 지점이 꽉 차서 지금은 사용할 수 없습니다."

const (
	ExpeditionJoinedMessage    = "%s님이 원정대에 참가하였습니다."
	ExpeditionLeftMessage      = "%s님이 원정대에서 탈퇴하였습니다."
	ExpeditionKickedMessage    = "원정대장이 귀하를 원정대 제재 대상에 등록하였습니다."
	ExpeditionAllowedMessage   = "원정대장이 귀하의 원정대 참여를 허가하였습니다."
	ExpeditionExpiredMessage   = "원정대 모집 시간이 종료되었습니다."
	ExpeditionDisbandedMessage = "원정대장이 자리를 비워 원정대가 해산되었습니다."
	ExpeditionSkippedMessage   = "%s님이 자리에 없어 원정대 대기 순서를 넘깁니다."
)
