package wz

type QuestRequirementKind string

const (
	QuestReqNPC             QuestRequirementKind = "npc"
	QuestReqLvMin           QuestRequirementKind = "lvmin"
	QuestReqLvMax           QuestRequirementKind = "lvmax"
	QuestReqClass           QuestRequirementKind = "job"
	QuestReqItem            QuestRequirementKind = "item"
	QuestReqMob             QuestRequirementKind = "mob"
	QuestReqQuest           QuestRequirementKind = "quest"
	QuestReqSkill           QuestRequirementKind = "skill"
	QuestReqPop             QuestRequirementKind = "pop"
	QuestReqInterval        QuestRequirementKind = "interval"
	QuestReqFieldEnter      QuestRequirementKind = "fieldEnter"
	QuestReqQuestComplete   QuestRequirementKind = "questComplete"
	QuestReqPet             QuestRequirementKind = "pet"
	QuestReqPetTamenessMin  QuestRequirementKind = "pettamenessmin"
	QuestReqMBMin           QuestRequirementKind = "mbmin"
	QuestReqMBCard          QuestRequirementKind = "mbcard"
	QuestReqSubClassFlags   QuestRequirementKind = "subJobFlags"
	QuestReqDayByDay        QuestRequirementKind = "dayByDay"
	QuestReqNormalAutoStart QuestRequirementKind = "normalAutoStart"
	QuestReqPartyQuestS     QuestRequirementKind = "partyQuest_S"
	QuestReqStartScript     QuestRequirementKind = "startscript"
	QuestReqEndScript       QuestRequirementKind = "endscript"
	QuestReqTimeStart       QuestRequirementKind = "start"
	QuestReqTimeEnd         QuestRequirementKind = "end"
	QuestReqInfo            QuestRequirementKind = "info"
	QuestReqInfoNumber      QuestRequirementKind = "infoNumber"
)

type QuestActionKind string

const (
	QuestActEXP        QuestActionKind = "exp"
	QuestActItem       QuestActionKind = "item"
	QuestActNextQuest  QuestActionKind = "nextQuest"
	QuestActMoney      QuestActionKind = "money"
	QuestActQuest      QuestActionKind = "quest"
	QuestActSkill      QuestActionKind = "skill"
	QuestActPop        QuestActionKind = "pop"
	QuestActBuffItemID QuestActionKind = "buffItemID"
	QuestActInfoNumber QuestActionKind = "infoNumber"
	QuestActSP         QuestActionKind = "sp"
	QuestActInfo       QuestActionKind = "info"
	QuestActNPCAct     QuestActionKind = "npcAct"
)

type Quest struct {
	ID         uint32
	Meta       QuestMeta
	Start      QuestPhase
	Complete   QuestPhase
	PartyRanks map[string][]PartyQuestRankCheck
}

type PartyQuestRankMode string

const (
	PartyQuestRankLess  PartyQuestRankMode = "less"
	PartyQuestRankMore  PartyQuestRankMode = "more"
	PartyQuestRankEqual PartyQuestRankMode = "equal"
)

type PartyQuestRankCheck struct {
	Mode     PartyQuestRankMode
	Property string
	Value    int
}

type QuestMeta struct {
	Name            string
	Parent          string
	Order           int
	Descriptions    map[int]string
	Area            int
	AutoStart       bool
	AutoPreComplete bool
	AutoComplete    bool
	AutoAccept      bool
	Repeatable      bool
	Blocked         bool
	ViewMedalItem   int
	SelectedSkillID int
	TimeLimit       int
	TimeLimit2      int
}

type QuestPhase struct {
	Requirements []QuestRequirement
	Actions      []QuestAction
}

type QuestRequirement struct {
	Kind        QuestRequirementKind
	IntValue    int
	StrValue    string
	InfoStrings []string
	Classes     []int
	PetIDs      []uint32
	Items       map[uint32]int
	Mobs        map[uint32]int
	Quests      map[uint32]QuestStatus
	Skills      map[uint32]int
}

type QuestAction struct {
	Kind              QuestActionKind
	IntValue          int
	StrValue          string
	ApplicableClasses []int
	Items             []QuestRewardItem
	Skills            []QuestRewardSkill
	Quests            map[uint32]QuestStatus
}

type QuestStatus uint8

const (
	QuestStatusNotStarted QuestStatus = 0
	QuestStatusStarted    QuestStatus = 1
	QuestStatusCompleted  QuestStatus = 2
)

type QuestRewardProp int

const (
	QuestRewardPropAlways QuestRewardProp = -2
	QuestRewardPropSelect QuestRewardProp = -1
)

func (p QuestRewardProp) IsAlways() bool {
	return p == QuestRewardPropAlways
}

func (p QuestRewardProp) IsSelection() bool {
	return p == QuestRewardPropSelect
}

func (p QuestRewardProp) IsWeightedRandom() bool {
	return p > 0
}

func (p QuestRewardProp) RandomWeight() int {
	if p <= 0 {
		return 0
	}
	return int(p)
}

type QuestRewardItem struct {
	ItemID     uint32
	Count      int
	Class      int
	ClassEx    int
	Gender     int
	Period     int
	Prop       QuestRewardProp
	DateExpire string
}

type QuestRewardSkill struct {
	SkillID     uint32
	SkillLevel  int
	MasterLevel int
	Classes     []int
}

func (q *Quest) RelevantMobs() map[uint32]int {
	if q == nil {
		return nil
	}
	out := make(map[uint32]int)
	for _, reqs := range [][]QuestRequirement{q.Start.Requirements, q.Complete.Requirements} {
		for _, req := range reqs {
			if req.Kind != QuestReqMob {
				continue
			}
			for mobID, count := range req.Mobs {
				out[mobID] = count
			}
		}
	}
	return out
}

func (q *Quest) OrderedMobIDs() []uint32 {
	if q == nil {
		return nil
	}
	out := make([]uint32, 0)
	seen := make(map[uint32]struct{})
	for _, reqs := range [][]QuestRequirement{q.Complete.Requirements, q.Start.Requirements} {
		for _, req := range reqs {
			if req.Kind != QuestReqMob {
				continue
			}
			for mobID := range req.Mobs {
				if _, ok := seen[mobID]; ok {
					continue
				}
				seen[mobID] = struct{}{}
				out = append(out, mobID)
			}
		}
	}
	return out
}

func (q *Quest) HasAutoStartMeta() bool {
	if q == nil {
		return false
	}
	return q.Meta.AutoStart || q.Meta.AutoAccept
}

func (q *Quest) IsPartyQuest() bool {
	return q != nil && len(q.PartyRanks) > 0
}

func (q *Quest) PartyRankChecks(rank string) []PartyQuestRankCheck {
	if q == nil || q.PartyRanks == nil {
		return nil
	}
	return q.PartyRanks[rank]
}
