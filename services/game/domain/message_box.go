package domain

import (
	"time"

	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

const (
	MessageBoxLifetime  = 10 * time.Minute
	messageBoxLines     = 3
	messageBoxLineLimit = 60
	messageBoxSpacingSq = 30000
)

type MessageBox struct {
	ObjectCore
	ItemID uint32
	Owner  string
	Text   string
}

func (box *MessageBox) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeMessageBox
}

func (box *MessageBox) Is(typ constant.ObjectType) bool {
	return box.GetObjectType().Has(typ)
}

func (box *MessageBox) SendSpawnSyncToViewer(viewer *Character) {
	viewer.Send(&response.SpawnMessageBox{
		OID:    box.OID,
		ItemID: box.ItemID,
		Text:   box.Text,
		Owner:  box.Owner,
		X:      box.Position.X,
		Y:      box.Position.Y,
	}, types.SEND_POLICY_ENCRYPT)
}

func (box *MessageBox) SendDestroySyncToViewer(viewer *Character) {
	viewer.Send(&response.RemoveMessageBox{OID: box.OID}, types.SEND_POLICY_ENCRYPT)
}
