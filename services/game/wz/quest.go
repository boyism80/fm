package wz

type QuestRequirementKind string

const (
	QuestReqNPC             QuestRequirementKind = "npc"
	QuestReqLvMin           QuestRequirementKind = "lvmin"
	QuestReqLvMax           QuestRequirementKind = "lvmax"
	QuestReqJob             QuestRequirementKind = "job"
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
	QuestReqSubJobFlags     QuestRequirementKind = "subJobFlags"
	QuestReqDayByDay        QuestRequirementKind = "dayByDay"
	QuestReqNormalAutoStart QuestRequirementKind = "normalAutoStart"
	QuestReqPartyQuestS     QuestRequirementKind = "partyQuest_S"
	QuestReqStartScript     QuestRequirementKind = "startscript"
	QuestReqEndScript       QuestRequirementKind = "endscript"
	QuestReqTimeStart       QuestRequirementKind = "start"
	QuestReqTimeEnd         QuestRequirementKind = "end"
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
	ID       uint32
	Meta     QuestMeta
	Start    QuestPhase
	Complete QuestPhase
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
	Blocked         bool
	ViewMedalItem   int
	SelectedSkillID int
}

type QuestPhase struct {
	Requirements []QuestRequirement
	Actions      []QuestAction
}

type QuestRequirement struct {
	Kind     QuestRequirementKind
	IntValue int
	StrValue string
	Jobs     []int
	PetIDs   []uint32
	Items    []QuestItemCount
	Mobs     []QuestMobCount
	Quests   []QuestStateRef
	Skills   []QuestSkillRef
}

type QuestAction struct {
	Kind           QuestActionKind
	IntValue       int
	StrValue       string
	ApplicableJobs []int
	Items          []QuestRewardItem
	Skills         []QuestRewardSkill
	Quests         []QuestStateRef
}

type QuestItemCount struct {
	ItemID uint32
	Count  int
}

type QuestMobCount struct {
	MobID uint32
	Count int
}

type QuestStateRef struct {
	QuestID uint32
	State   int
}

type QuestSkillRef struct {
	SkillID uint32
	Acquire int
}

type QuestRewardItem struct {
	ItemID     uint32
	Count      int
	Job        int
	JobEx      int
	Gender     int
	Period     int
	Prop       int
	DateExpire string
}

type QuestRewardSkill struct {
	SkillID     uint32
	SkillLevel  int
	MasterLevel int
	Jobs        []int
}

func (q *Quest) RelevantMobs() map[uint32]int {
	if q == nil {
		return nil
	}
	out := make(map[uint32]int)
	q.collectRelevantMobs(out, q.Start.Requirements)
	q.collectRelevantMobs(out, q.Complete.Requirements)
	return out
}

func (q *Quest) OrderedMobIDs() []uint32 {
	if q == nil {
		return nil
	}
	out := make([]uint32, 0)
	seen := make(map[uint32]struct{})
	q.appendOrderedMobIDs(&out, seen, q.Complete.Requirements)
	q.appendOrderedMobIDs(&out, seen, q.Start.Requirements)
	return out
}

func (q *Quest) appendOrderedMobIDs(out *[]uint32, seen map[uint32]struct{}, reqs []QuestRequirement) {
	for _, req := range reqs {
		if req.Kind != QuestReqMob {
			continue
		}
		for _, mob := range req.Mobs {
			if _, ok := seen[mob.MobID]; ok {
				continue
			}
			seen[mob.MobID] = struct{}{}
			*out = append(*out, mob.MobID)
		}
	}
}

func (q *Quest) collectRelevantMobs(out map[uint32]int, reqs []QuestRequirement) {
	for _, req := range reqs {
		if req.Kind != QuestReqMob {
			continue
		}
		for _, mob := range req.Mobs {
			out[mob.MobID] = mob.Count
		}
	}
}
