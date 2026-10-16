package game

import (
	"time"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

const (
	MessageBoxLifetime  = 10 * time.Minute
	messageBoxLines     = 3
	messageBoxLineLimit = 60
	messageBoxSpacingSq = 30000
)

type MessageBox struct {
	objectCore
	itemID uint32
	owner  string
	text   string
}

func NewMessageBox(owner *Character, itemID uint32, text string) *MessageBox {
	m := owner.GetMap()
	if m == nil || !owner.IsAlive() || m.wz.MessageBoxBlocked || !ClientText(text).ValidLines(messageBoxLines, messageBoxLineLimit) {
		return nil
	}
	box := &MessageBox{itemID: itemID, owner: owner.GetName(), text: text}
	box.position = owner.position
	return box
}

func (box *MessageBox) Place(owner *Character) bool {
	if err := owner.GetMap().AddMessageBox(box); err != nil {
		owner.listener.OnMessageBoxFailed(owner)
		return false
	}
	return true
}

func (box *MessageBox) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeMessageBox
}

func (box *MessageBox) Is(typ constant.ObjectType) bool {
	return box.GetObjectType().Has(typ)
}

func (box *MessageBox) SendSpawnSyncToViewer(viewer *Character) {
	viewer.Send(&response.SpawnMessageBox{
		OID:    box.oid,
		ItemID: box.itemID,
		Text:   box.text,
		Owner:  box.owner,
		X:      box.position.X,
		Y:      box.position.Y,
	}, types.SEND_POLICY_ENCRYPT)
}

func (box *MessageBox) SendDestroySyncToViewer(viewer *Character) {
	viewer.Send(&response.RemoveMessageBox{OID: box.oid}, types.SEND_POLICY_ENCRYPT)
}
