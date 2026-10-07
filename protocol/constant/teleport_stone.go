package constant

type TeleportStoneAction uint8

const (
	TeleportStoneActionRemove   TeleportStoneAction = 0
	TeleportStoneActionRegister TeleportStoneAction = 1
)

type TeleportStoneResult uint8

const (
	TeleportStoneResultList           TeleportStoneResult = 3
	TeleportStoneResultCannotGo       TeleportStoneResult = 5
	TeleportStoneResultNotFound       TeleportStoneResult = 6
	TeleportStoneResultCurrentMap     TeleportStoneResult = 9
	TeleportStoneResultCannotRegister TeleportStoneResult = 10
)
