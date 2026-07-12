package wz

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
	Requirements QuestRequirements
	Actions      QuestActions
}

type QuestRequirements struct {
	NPC               uint32
	LevelMin          int
	LevelMax          int
	Level             int
	Job               []int
	Item              map[uint32]int
	Mob               map[uint32]int
	Quest             map[uint32]QuestStatus
	Skill             map[uint32]int
	Pop               int
	Interval          int
	HasInterval       bool
	FieldEnter        int
	QuestComplete     int
	Pet               []uint32
	PetTamenessMin    int
	MBMin             int
	MBCard            map[uint32]int
	SubJobFlags       int
	DayByDay          bool
	NormalAutoStart   bool
	PartyQuestS       int
	StartScript       string
	EndScript         string
	Start             string
	End               string
	Info              []string
	InfoNumber        int
	WorldMin          string
	WorldMax          string
	EndMeso           int
	EquipAllNeed      int
	EquipSelectNeed   int
	Premium           bool
	Buff              string
	ExceptBuff        string
	TamingMobLevelMin int
}

type QuestActions struct {
	Item        []QuestActionItem
	Exp         int
	Money       int
	Pop         int
	NextQuest   uint32
	BuffItemID  uint32
	Info        string
	NPCAct      string
	NPC         int
	Quests      map[uint32]QuestStatus
	Skills      []QuestActionSkill
	SkillJobs   []int
	SP          int
	SPJobs      []int
	InfoNumber  uint32
	PetTameness int
	PetSpeed    int
	Map         int
	Job         int
	LvMin       int
	LvMax       int
	FieldEnter  int
	Interval    int
	Message     string
	Start       string
	End         string
	Ask         int
	Stop        int
	Say         map[string]string
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

type QuestActionItem struct {
	ItemID     uint32
	Count      int
	Class      int
	ClassEx    int
	Gender     int
	Period     int
	Prop       QuestRewardProp
	DateExpire string
}

type QuestActionSkill struct {
	SkillID     uint32
	SkillLevel  int
	MasterLevel int
	Classes     []int
}

type QuestRewardItem = QuestActionItem
type QuestRewardSkill = QuestActionSkill

func (q *Quest) RelevantMobs() map[uint32]int {
	if q == nil {
		return nil
	}
	out := make(map[uint32]int)
	for mobID, count := range q.Start.Requirements.Mob {
		out[mobID] = count
	}
	for mobID, count := range q.Complete.Requirements.Mob {
		out[mobID] = count
	}
	return out
}

func (q *Quest) OrderedMobIDs() []uint32 {
	if q == nil {
		return nil
	}
	out := make([]uint32, 0)
	seen := make(map[uint32]struct{})
	for _, mobs := range []map[uint32]int{q.Complete.Requirements.Mob, q.Start.Requirements.Mob} {
		for mobID := range mobs {
			if _, ok := seen[mobID]; ok {
				continue
			}
			seen[mobID] = struct{}{}
			out = append(out, mobID)
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
