package server

import (
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnQuestStarted(ch *game.Character, qp *game.Quest, npcID uint32) {
	status, ok := qp.BuildS2C()
	if !ok {
		return
	}
	ch.Send(&response.UpdateQuest{QuestStatus: status}, types.SEND_POLICY_ENCRYPT)
	if npcID != 0 {
		ch.Send(&response.UpdateQuestNPC{
			Progress:    8,
			QuestID:     status.QuestID,
			NPCID:       npcID,
			NextQuestID: 0,
		}, types.SEND_POLICY_ENCRYPT)
	}
	if limit := qp.TimeLimit(); limit > 0 {
		l.OnClock(ch, int32(limit))
	}
	l.OnUnlockAction(ch)
}

func (l *CharacterListenerImpl) OnQuestCompleted(ch *game.Character, qp *game.Quest, npcID uint32, nextQuestID uint32) {
	status, ok := qp.BuildS2C()
	if !ok {
		return
	}
	ch.Send(&response.UpdateQuest{QuestStatus: status}, types.SEND_POLICY_ENCRYPT)
	if npcID != 0 {
		ch.Send(&response.UpdateQuestNPC{
			Progress:    8,
			QuestID:     status.QuestID,
			NPCID:       npcID,
			NextQuestID: nextQuestID,
		}, types.SEND_POLICY_ENCRYPT)
	}
	l.OnUnlockAction(ch)
}

func (l *CharacterListenerImpl) OnQuestForfeited(ch *game.Character, qp *game.Quest) {
	status, ok := qp.BuildS2C()
	if !ok {
		return
	}
	ch.Send(&response.UpdateQuest{QuestStatus: status}, types.SEND_POLICY_ENCRYPT)
	l.OnUnlockAction(ch)
}

func (l *CharacterListenerImpl) OnShowQuestCompletion(ch *game.Character, questID uint32) {
	if questID == 0 || questID > 0xFFFF {
		return
	}
	ch.Send(&response.ShowQuestCompletion{
		QuestID: uint16(questID),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnQuestRecordChanged(ch *game.Character, qp *game.Quest) {
	status, ok := qp.BuildS2C()
	if !ok {
		return
	}
	status.Status = response.QuestWireStatusStarted
	ch.Send(&response.UpdateQuest{QuestStatus: status}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnQuestProgress(ch *game.Character, qp *game.Quest) {
	status, ok := qp.BuildS2C()
	if !ok {
		return
	}
	l.OnQuestRecordChanged(ch, qp)
	if qp.CanComplete(ch, game.QuestPhaseOpts{}) == nil {
		ch.Send(&response.ShowQuestCompletion{QuestID: status.QuestID}, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *CharacterListenerImpl) OnQuestRecordExChanged(ch *game.Character, qp *game.Quest) {
	status, ok := qp.BuildS2C()
	if !ok {
		return
	}
	ch.Send(&response.UpdateQuestRecordEx{
		QuestID: status.QuestID,
		Data:    qp.RecordExWire(),
	}, types.SEND_POLICY_ENCRYPT)
}
